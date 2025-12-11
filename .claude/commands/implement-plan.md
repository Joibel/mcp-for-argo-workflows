# Implement Plan from Linear

You are orchestrating the implementation of the mcp-for-argo-workflows project based on the Linear project plan.

## Step 1: Fetch Current State from Linear

First, retrieve the current project status:

1. Use `mcp__linear-server__list_issues` with `project: "mcp-for-argo-workflows"` to get all issues
2. Identify issues by status: Backlog, Todo, In Progress, Done
3. Check issue dependencies (noted in each issue's description)

## Step 2: Identify Ready Tasks

Analyze the issues to find tasks that are ready to implement:

1. **Setup tasks** (PIP-5 through PIP-14) - Foundation work, check dependencies
2. **MCP tools** (PIP-15 through PIP-49) - Require setup tasks complete
3. **Testing** (PIP-50, PIP-51) - Require tools implemented
4. **Documentation** (PIP-9, PIP-52, PIP-53) - Can be done incrementally

A task is "ready" when:
- Status is Backlog or Todo
- All dependencies (listed in issue description) are completed

## Step 3: Present Ready Tasks

Present the list of ready tasks to the user with:
- Issue identifier (e.g., PIP-10)
- Title
- Brief description of what it involves
- Dependencies (and their status)

## Step 4: Execute Implementation via /implement-issue

For each ready task, use the `/implement-issue` command to implement it:

```
/implement-issue PIP-XX
```

**IMPORTANT**: Always use `/implement-issue` to implement individual tasks. This ensures consistent implementation workflow across all tasks.

Wait for each task to complete before starting the next one.

## Step 5: Report Progress

After each implementation cycle:

1. Summarize what was completed
2. List any blockers or issues encountered
3. Identify next tasks ready for implementation (dependencies now met)
4. Ask user if they want to continue with the next task

## Implementation Guidelines

- **One task at a time** - Complete and verify before moving to next
- **Always use /implement-issue** - Never implement tasks directly in this command
- **Ask for clarification** - If requirements are ambiguous, ask before implementing
- **Report blockers** - If a task cannot be completed, report why and suggest alternatives

## Starting Point

If this is the first implementation session, the typical order is:
1. PIP-5: Initialize Go module and directory structure
2. PIP-6: Create Makefile
3. PIP-7: Configure golangci-lint
4. PIP-8: Set up GitHub Actions CI
5. PIP-9: Add basic README

These have no dependencies and enable all subsequent work.

---

Begin by fetching the current Linear project state and identifying ready tasks.
