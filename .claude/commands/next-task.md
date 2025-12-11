# Get Next Task

Identify and suggest the next task to implement from the Linear project plan.

## Step 1: Fetch Project State

1. Get all issues: `mcp__linear-server__list_issues` with `project: "mcp-for-argo-workflows"`
2. Filter by status to find Backlog/Todo items
3. Check what's currently In Progress

## Step 2: Analyze Dependencies

For each Backlog/Todo issue:
1. Extract dependencies from issue description
2. Check if all dependencies are Done
3. Mark issue as "ready" if dependencies met

## Step 3: Prioritize Ready Tasks

Apply priority order:

1. **Setup tasks first** (PIP-5 → PIP-14)
   - These enable all other work
   - Order: module init → Makefile → linting → CI → server skeleton → transports → client

2. **Core MCP tools** (PIP-15 → PIP-22)
   - submit_workflow, list_workflows, get_workflow, delete_workflow
   - logs, watch, wait
   - lint_workflow (validation)

3. **Control tools** (PIP-23 → PIP-28)
   - retry, resubmit, suspend, resume, stop, terminate

4. **Template tools** (PIP-29 → PIP-42)
   - WorkflowTemplates, ClusterWorkflowTemplates, CronWorkflows

5. **Archive tools** (PIP-43 → PIP-47)
   - Require Argo Server mode

6. **Node tools** (PIP-48 → PIP-49)
   - Advanced operations

7. **Testing** (PIP-50 → PIP-51)
   - Can be done alongside tool implementation

8. **Documentation** (PIP-52 → PIP-53)
   - Final polish

## Step 4: Check Current Progress

1. What's already implemented locally?
2. Are there any In Progress items that should be completed first?
3. Any failing tests blocking progress?

## Step 5: Recommend Next Task

Output format:

```
## Next Recommended Task

**Issue:** PIP-X - Issue Title
**Status:** Backlog
**Dependencies:** All met ✓
**Specialist:** agent-name
**Estimated Scope:** Small/Medium/Large

### Why This Task?
Brief explanation of why this is the best next task.

### Quick Start
To implement this task, run:
`/implement-issue PIP-X`

### Alternative Tasks
If you'd prefer something different:
- PIP-Y - Also ready, different scope
- PIP-Z - Also ready, different area
```

---

Begin by fetching the current Linear project state.
