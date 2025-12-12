# Review Project Progress

Review the current state of the mcp-for-argo-workflows project against the Linear plan.

## Step 1: Gather Current State

1. **Fetch Linear issues** using `mcp__linear-server__list_issues` with `project: "mcp-for-argo-workflows"`
2. **Check local codebase** - What files exist? What's implemented?
3. **Run verification** - Execute `make lint` and `make test` if Makefile exists

## Step 2: Identify Discrepancies

Report on:
1. **Issues marked Done in Linear but not implemented locally**
2. **Code implemented locally but issues still in Backlog**
3. **Failing tests or lint errors**
4. **Missing dependencies or blockers**

## Step 3: Generate Report

Produce a summary with:

```
## Progress Report: mcp-for-argo-workflows

### Overall Status
- Setup: X/10 complete
- MCP Tools: X/35 complete
- Testing: X/2 complete
- Documentation: X/3 complete

### Recently Completed
- [PIP-X] Issue title

### In Progress
- [PIP-X] Issue title - status notes

### Ready for Implementation
- [PIP-X] Issue title (dependencies met)

### Blocked
- [PIP-X] Issue title - blocked by [PIP-Y]

### Recommendations
1. Next priority task
2. Any issues needing attention
```

## Step 4: Update Linear (Optional)

If discrepancies found, offer to:
1. Update issue statuses to match reality
2. Add progress comments to in-progress issues
3. Flag any blockers or concerns

---

Begin by fetching the Linear project state.
