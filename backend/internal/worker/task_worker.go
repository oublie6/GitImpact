package worker

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"gitimpact/backend/internal/analyzer"
	"gitimpact/backend/internal/config"
	"gitimpact/backend/internal/model"
	"gitimpact/backend/internal/repository"
)

type TaskWorker struct {
	cfg        *config.AppConfig
	analyzer   analyzer.Analyzer
	taskRepo   *repository.TaskRepository
	repoRepo   *repository.RepoRepository
	reportRepo *repository.ReportRepository
}

func NewTaskWorker(cfg *config.AppConfig, a analyzer.Analyzer, tr *repository.TaskRepository, rr *repository.RepoRepository, rp *repository.ReportRepository) *TaskWorker {
	return &TaskWorker{cfg: cfg, analyzer: a, taskRepo: tr, repoRepo: rr, reportRepo: rp}
}

func (w *TaskWorker) Enqueue(taskID uint) { go w.process(context.Background(), taskID) }

func (w *TaskWorker) process(ctx context.Context, taskID uint) {
	task, err := w.taskRepo.GetByID(taskID)
	if err != nil {
		return
	}
	task.Status = model.TaskStatusRunning
	_ = w.taskRepo.Update(task)
	_ = w.taskRepo.AddLog(&model.TaskLog{TaskID: taskID, Level: "INFO", Message: "task started"})

	oldRepo, err := w.repoRepo.GetByID(task.OldRepoID)
	if err != nil {
		w.fail(task, err)
		return
	}
	newRepo, err := w.repoRepo.GetByID(task.NewRepoID)
	if err != nil {
		w.fail(task, err)
		return
	}

	workDir := filepath.Join(w.cfg.Workdir.Artifacts, fmt.Sprintf("task_%d_%d", task.ID, time.Now().Unix()))
	_ = os.MkdirAll(workDir, 0o755)
	if err := w.prepareRepo(oldRepo, task.OldRef); err != nil {
		w.fail(task, err)
		return
	}
	if err := w.prepareRepo(newRepo, task.NewRef); err != nil {
		w.fail(task, err)
		return
	}
	if err := w.writeMaterials(workDir, oldRepo.LocalCacheDir, task.OldRef, newRepo.LocalCacheDir, task.NewRef, task.CustomFocus); err != nil {
		w.fail(task, err)
		return
	}

	mdOut, mdErr, mdExecErr := w.analyzer.RunMarkdownReport(ctx, workDir)
	jsonOut, jsonErr, _ := w.analyzer.RunStructuredReport(ctx, workDir)
	report := &model.AnalysisReport{TaskID: taskID, MarkdownReport: mdOut, StructuredReport: jsonOut, RawStdout: strings.TrimSpace(mdOut + "\n" + jsonOut), RawStderr: strings.TrimSpace(mdErr + "\n" + jsonErr)}
	_ = w.reportRepo.Upsert(taskID, report)
	if mdExecErr != nil {
		w.fail(task, mdExecErr)
		return
	}
	task.Status = model.TaskStatusSuccess
	task.ErrorMessage = ""
	_ = w.taskRepo.Update(task)
	_ = w.taskRepo.AddLog(&model.TaskLog{TaskID: taskID, Level: "INFO", Message: "task completed"})
}

func (w *TaskWorker) fail(task *model.AnalysisTask, err error) {
	task.Status = model.TaskStatusFailed
	task.ErrorMessage = err.Error()
	_ = w.taskRepo.Update(task)
	_ = w.taskRepo.AddLog(&model.TaskLog{TaskID: task.ID, Level: "ERROR", Message: err.Error()})
}

func (w *TaskWorker) prepareRepo(repo *model.Repository, ref string) error {
	if _, err := os.Stat(filepath.Join(repo.LocalCacheDir, ".git")); err != nil {
		if out, err := exec.Command("git", "clone", repo.RepoURL, repo.LocalCacheDir).CombinedOutput(); err != nil {
			return fmt.Errorf("clone repo failed %s: %w", string(out), err)
		}
	}
	if out, err := exec.Command("git", "-C", repo.LocalCacheDir, "fetch", "--all", "--prune").CombinedOutput(); err != nil {
		return fmt.Errorf("fetch failed %s: %w", string(out), err)
	}
	if out, err := exec.Command("git", "-C", repo.LocalCacheDir, "checkout", ref).CombinedOutput(); err != nil {
		return fmt.Errorf("checkout failed %s: %w", string(out), err)
	}
	return nil
}

func (w *TaskWorker) writeMaterials(workDir, oldDir, oldRef, newDir, newRef, focus string) error {
	write := func(name, content string) error {
		return os.WriteFile(filepath.Join(workDir, name), []byte(content), 0o644)
	}
	if out, err := exec.Command("git", "-C", newDir, "diff", "--name-only", oldRef, newRef).CombinedOutput(); err == nil {
		_ = write("changed_files.txt", string(out))
	}
	if out, err := exec.Command("git", "-C", newDir, "diff", oldRef, newRef).CombinedOutput(); err == nil {
		_ = write("diff.patch", string(out))
	}
	if out, err := exec.Command("git", "-C", newDir, "log", "--oneline", oldRef+".."+newRef).CombinedOutput(); err == nil {
		_ = write("commit_log.txt", string(out))
	}
	manifest := fmt.Sprintf("# Repo Manifest\n- old_ref: %s\n- new_ref: %s\n- old_repo: %s\n- new_repo: %s\n", oldRef, newRef, oldDir, newDir)
	_ = write("repo_manifest.md", manifest)
	prompt := "请结合 changed_files.txt、diff.patch、commit_log.txt、repo_manifest.md 生成 Markdown 影响分析报告。\n关注点: " + focus
	promptJSON := "请输出结构化 JSON，字段包括 summary/changed_modules/impacted_interfaces/impacted_configs/impacted_scripts/impacted_tests/risks/backward_compatibility/deployment_risks/rollback_risks/verification_suggestions/confidence/raw_notes。\n关注点: " + focus
	_ = write("analysis_prompt.md", prompt)
	_ = write("analysis_prompt_json.md", promptJSON)
	return nil
}
