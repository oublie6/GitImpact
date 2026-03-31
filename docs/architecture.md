# 架构说明
- Monorepo: `backend/` + `frontend/` + `sql/` + `scripts/` + `docs/` + `examples/`。
- 后端分层：router / handler / service / repository / model / analyzer / worker / middleware / config。
- worker 异步执行任务，主线程只负责入库与入队。
