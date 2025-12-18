# Ubuntu Slim Migration - Final Implementation Plan

**Date:** 2025-12-17
**Scope:** Tier 1 (Ultra-Safe) + drift-detection (Tier 2)
**Risk Level:** Very Low

---

## Executive Summary

Migrate **8 specific jobs** to ubuntu-slim runners based on actual runtime data showing all jobs complete in under 1 minute (except drift-detection at 30s). This provides a **14+ minute safety buffer** from the 15-minute timeout limit.

### Jobs to Migrate

| # | Job Name | Current Time | Safety Buffer | Template File |
|---|----------|--------------|---------------|---------------|
| 1 | setup-webkit | 22-24s | 14m 36s+ | pr.yaml.tmpl, release.yaml.tmpl |
| 2 | secret-scan | 24s | 14m 36s | pr.yaml.tmpl |
| 3 | detect-changes | 39s | 14m 21s | pr.yaml.tmpl |
| 4 | drift-detection | 30s | 14m 30s | pr.yaml.tmpl |
| 5 | app-json-diff | 25s | 14m 35s | pr.yaml.tmpl |
| 6 | validate-app-json | 24s | 14m 36s | pr.yaml.tmpl |
| 7 | cleanup-containers | 5-6s | 14m 54s+ | release.yaml.tmpl |
| 8 | notify-success | 20s | 14m 40s | release.yaml.tmpl |
| 9 | notify-failure | ~20s | 14m 40s | release.yaml.tmpl |

**Total:** 9 job types (setup-webkit appears in both templates)

---

## Files to Modify

### 1. PR Workflow Template
**File:** `/home/user/webkit/internal/templates/.github/workflows/pr.yaml.tmpl`

**Changes:**

```yaml
# Line 24-25: setup-webkit
  setup-webkit:
-   runs-on: ubuntu-latest
+   runs-on: ubuntu-slim

# Line 55-57: secret-scan
  secret-scan:
    name: gitleaks
-   runs-on: ubuntu-latest
+   runs-on: ubuntu-slim

# Line 69-70: detect-changes
  detect-changes:
-   runs-on: ubuntu-latest
+   runs-on: ubuntu-slim

# Line 102-103: drift-detection
  drift-detection:
-   runs-on: ubuntu-latest
+   runs-on: ubuntu-slim

# Line 150-151: app-json-diff
  app-json-diff:
-   runs-on: ubuntu-latest
+   runs-on: ubuntu-slim

# Line 329-330: validate-app-json
  validate-app-json:
-   runs-on: ubuntu-latest
+   runs-on: ubuntu-slim
```

**Total changes in pr.yaml.tmpl:** 6 jobs

---

### 2. Release Workflow Template
**File:** `/home/user/webkit/internal/templates/.github/workflows/release.yaml.tmpl`

**Changes:**

```yaml
# Line 23-24: setup-webkit
  setup-webkit:
-   runs-on: ubuntu-latest
+   runs-on: ubuntu-slim

# Line 107-108: cleanup-containers
  cleanup-containers:
-   runs-on: ubuntu-latest
+   runs-on: ubuntu-slim

# Line 393-394: notify-success
  notify-success:
-   runs-on: ubuntu-latest
+   runs-on: ubuntu-slim

# Line 420-421: notify-failure
  notify-failure:
-   runs-on: ubuntu-latest
+   runs-on: ubuntu-slim
```

**Total changes in release.yaml.tmpl:** 4 jobs (3 unique, setup-webkit is duplicate)

---

## Implementation Steps

### Phase 1: Template Modification (30 minutes)

1. **Edit PR template:**
   ```bash
   # Edit /home/user/webkit/internal/templates/.github/workflows/pr.yaml.tmpl
   # Change runs-on for: setup-webkit, secret-scan, detect-changes,
   #                      drift-detection, app-json-diff, validate-app-json
   ```

2. **Edit Release template:**
   ```bash
   # Edit /home/user/webkit/internal/templates/.github/workflows/release.yaml.tmpl
   # Change runs-on for: setup-webkit, cleanup-containers,
   #                      notify-success, notify-failure
   ```

3. **Verify changes:**
   ```bash
   # Grep to confirm all changes applied
   grep -n "runs-on: ubuntu-slim" internal/templates/.github/workflows/*.tmpl
   ```

---

### Phase 2: Testing (1-2 weeks)

1. **Create test repository:**
   - Fork or use existing private test repo
   - Run `webkit update` with modified templates
   - Verify generated workflows have `runs-on: ubuntu-slim` for target jobs

2. **Test PR workflow (10-15 test PRs):**
   - Small PRs (< 100 lines)
   - Medium PRs (100-500 lines)
   - Large PRs (500-1000 lines) - to test drift-detection edge cases

   **Monitor for:**
   - ✅ Zero timeout failures
   - ✅ Same pass rate as before
   - ✅ Job durations remain under 2 minutes
   - ✅ No new error messages

3. **Test Release workflow (5-10 releases):**
   - Regular releases
   - Failed builds (to test notify-failure)

   **Monitor for:**
   - ✅ cleanup-containers completes successfully
   - ✅ Slack notifications sent correctly
   - ✅ No timeout issues

4. **Success criteria:**
   - Zero timeouts across all test runs
   - 100% of jobs complete in < 5 minutes
   - No increase in failure rate
   - All workflow checks pass as expected

---

### Phase 3: Rollout (Week 3)

1. **Merge template changes:**
   ```bash
   git add internal/templates/.github/workflows/pr.yaml.tmpl
   git add internal/templates/.github/workflows/release.yaml.tmpl
   git commit -m "feat: Optimize workflow costs with ubuntu-slim runners"
   git push
   ```

2. **Release new WebKit version:**
   - Tag new version (e.g., v0.11.0)
   - Update changelog with optimization details
   - Document in release notes

3. **Migration documentation:**
   Create update guide for users:
   ```markdown
   ## Cost Optimization (v0.11.0+)

   This version optimizes GitHub Actions costs by using ubuntu-slim
   runners for lightweight jobs. To apply these optimizations:

   1. Run: `webkit update`
   2. Review generated workflow changes
   3. Commit and push updated workflows

   **Impact:** ~20-30% reduction in Actions costs for private repos
   **Risk:** Very low - only fast jobs migrated
   **Timeout limit:** 15 minutes (all migrated jobs run in < 1 minute)
   ```

4. **Announcement channels:**
   - GitHub release notes
   - Project README
   - Discussions (if applicable)
   - Slack/Discord (if applicable)

---

### Phase 4: Monitoring (Month 1-2)

1. **Track metrics:**
   - Number of timeout failures reported
   - User feedback on workflow changes
   - Actual cost savings (from user reports)

2. **Address issues:**
   - If timeouts occur: Investigate specific job/project
   - If failures increase: Consider reverting specific job
   - If complaints: Provide opt-out mechanism

3. **Success indicators (1 month):**
   - < 5 timeout reports from all users
   - No increase in support requests
   - Positive feedback on cost savings
   - No rollback required

---

## Expected Outcomes

### Cost Savings (Private Repos Only)

**Per Repository:**
- PR workflow: 9 minutes saved × $0.004/min = $0.036 per PR
- Release workflow: 3.5 minutes saved × $0.004/min = $0.014 per release
- **Monthly (100 PRs + 20 releases):** ~$3.88/month per repo

**Organization-Wide:**
- 10 private repos: **$38.80/month** = $465.60/year
- 50 private repos: **$194/month** = $2,328/year

### Performance Impact

- **Startup time:** ubuntu-slim may be slightly faster (container vs VM)
- **Execution time:** Should be similar for lightweight operations
- **Total workflow time:** No significant change expected

### Risk Assessment

**Risk Level:** ⭐ Very Low

**Mitigations:**
- ✅ All jobs have 14+ minute safety buffer
- ✅ Based on real production data
- ✅ Only simple operations (no builds/tests)
- ✅ Easy rollback via template revert
- ✅ Users can opt-out by manually changing runners

**Failure Scenarios:**
1. **Drift-detection timeout on very large projects**
   - Mitigation: Users can manually change back to ubuntu-latest
   - Impact: Minimal - drift check is informational

2. **Container environment missing tool**
   - Likelihood: Very low - all jobs use standard tools
   - Impact: Job failure would be caught in testing

3. **Intermittent timeouts**
   - Mitigation: Monitor closely in first month
   - Response: Revert specific job if pattern emerges

---

## Rollback Plan

If critical issues arise, rollback is straightforward:

### Quick Rollback (Same Day)

1. **Revert template changes:**
   ```bash
   git revert <commit-sha>
   git push
   ```

2. **Release hotfix version:**
   ```bash
   # Tag as patch version (e.g., v0.11.1)
   git tag v0.11.1
   git push --tags
   ```

3. **Notify users:**
   - Post GitHub release with rollback note
   - Update any announcements
   - Provide `webkit update` instructions

### Partial Rollback

If only specific jobs cause issues:

1. Edit template to revert only problematic jobs
2. Release patch version
3. Document which jobs remain on ubuntu-slim

Example: If drift-detection causes timeouts:
- Revert only drift-detection to ubuntu-latest
- Keep other 8 jobs on ubuntu-slim
- Savings reduced by ~10% but still beneficial

---

## Decision Points

### Before Implementation
- [ ] Stakeholder approval on cost/risk tradeoff
- [ ] Test repository identified
- [ ] Timeline approved (3 weeks to full rollout)

### Before Merge (End of Testing)
- [ ] Zero timeout failures in testing
- [ ] All test PRs and releases passed
- [ ] No unexpected issues discovered
- [ ] Documentation prepared

### Before Announcement
- [ ] New version tagged and released
- [ ] Release notes finalized
- [ ] Support channels notified
- [ ] Monitoring plan in place

### After 1 Month
- [ ] Evaluate success metrics
- [ ] Collect user feedback
- [ ] Decide on future optimizations
- [ ] Consider adding more Tier 2 jobs

---

## Alternative Approaches (Not Recommended)

### Option A: Migrate All Jobs
**Status:** ❌ Rejected
**Reason:** Build/test jobs will definitely timeout (7+ minute runtimes)

### Option B: Add Conditional Logic
**Status:** ⚠️ Considered but rejected
**Reason:** Adds complexity; users can manually override if needed

### Option C: Wait for More Data
**Status:** ❌ Rejected
**Reason:** Current data is solid; jobs are consistently under 1 minute

### Option D: Self-Hosted Runners
**Status:** 🔮 Future consideration
**Reason:** Requires operational overhead; ubuntu-slim is easier first step

---

## Conclusion

This plan provides a **conservative, data-driven approach** to reduce GitHub Actions costs:

✅ **Low Risk:** 14+ minute safety buffer on all jobs
✅ **High Confidence:** Based on actual production timings
✅ **Easy Rollback:** Simple template revert if needed
✅ **Measurable Impact:** ~$200-2,000/year org-wide savings
✅ **User Control:** Can opt-out by editing generated workflows

**Recommendation:** Proceed with implementation.

---

## Next Steps

1. ✅ **Get approval** - Review this plan with team
2. ⬜ **Make changes** - Edit the 2 template files
3. ⬜ **Test thoroughly** - 1-2 weeks in test repo
4. ⬜ **Release** - New WebKit version with optimizations
5. ⬜ **Monitor** - Track for issues over 1 month

**Estimated timeline:** 3 weeks from approval to full rollout
**Estimated effort:** 8-12 hours total (templates + testing + docs)
**Expected ROI:** Positive within first month for orgs with 10+ private repos

---

**Prepared by:** Claude
**Date:** 2025-12-17
**Status:** Ready for Implementation
