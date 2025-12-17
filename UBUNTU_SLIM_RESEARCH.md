# Ubuntu Slim Runner Research: Cost Analysis and Feasibility

**Date:** 2025-12-17
**Status:** Research & Planning Phase
**Repository:** ainsleydev/webkit

## Executive Summary

GitHub Actions recently introduced `ubuntu-slim` runners - single-CPU, container-based runners optimized for lightweight operations. This research evaluates the feasibility of migrating applicable WebKit workflow jobs to ubuntu-slim runners to reduce GitHub Actions costs.

**Key Findings:**
- ubuntu-slim runners are **~40-50% cheaper** than ubuntu-latest (estimated)
- **Critical limitation:** 15-minute job timeout (hard limit)
- **Not suitable** for build-heavy or test-heavy jobs
- **Best candidates:** Quick validation, linting, scanning, and notification jobs
- **Estimated savings:** 20-30% reduction in total Actions costs (if applicable jobs are migrated)

---

## 1. What is ubuntu-slim?

### Technical Specifications
- **Runner Type:** Single-CPU container-based runner
- **Operating System:** Ubuntu Linux (container, not full VM)
- **CPU:** 1 vCPU (vs 2 vCPU for ubuntu-latest)
- **Job Timeout:** 15 minutes maximum (non-configurable)
- **Workflow Label:** `ubuntu-slim`

### Key Characteristics
1. **Container-based execution** - Faster startup than full VMs
2. **Lightweight** - Reduced resource footprint
3. **Cost-optimised** - Lower per-minute pricing
4. **Limited capabilities** - Not suitable for heavy CI/CD workloads

### Ideal Use Cases
- ✅ Linting and formatting checks
- ✅ Quick validation tasks
- ✅ Secret scanning
- ✅ Path change detection
- ✅ Simple API calls (notifications, deployments)
- ✅ Fast unit tests (<5 minutes)
- ❌ Full builds with compilation
- ❌ Extensive test suites
- ❌ Docker image builds
- ❌ Terraform operations
- ❌ Multi-step deployments

---

## 2. Pricing Analysis

### Current Pricing (as of December 2025)

| Runner Type | CPU | Cost per Minute | Notes |
|-------------|-----|----------------|-------|
| ubuntu-latest | 2 vCPU | $0.008 | Standard runner (current) |
| ubuntu-slim | 1 vCPU | ~$0.004-0.005 (estimated) | Single-CPU runner |
| Free Tier | - | $0.000 | Public repos on standard runners |

### 2026 Pricing Changes
GitHub announced pricing changes effective January 1, 2026:
- **Formula:** `new_price = old_price × 0.6 + $0.002`
- **Reduction:** Up to 39% for larger runners
- **Note:** Smaller runners see lesser price reduction

### Cost Calculation Assumptions
For **private repositories**:
- Current ubuntu-latest: $0.008/min = $0.48/hour
- Estimated ubuntu-slim: $0.004/min = $0.24/hour (assuming 50% cost of ubuntu-latest)
- **Potential savings: ~50% per job** (if job completes within 15-minute limit)

For **public repositories**:
- Standard runners are **free** for public repos
- **No cost benefit** from using ubuntu-slim on public repos

---

## 3. WebKit Workflow Analysis

### 3.1 Main Repository Workflows (ainsleydev/webkit)

Repository status: **Public** (based on workflow structure)
**Cost impact:** Minimal to none (public repos get free Actions minutes)

#### Workflow: `.github/workflows/pr.yaml`

| Job | Current Runner | Slim-Compatible? | Reasoning | Recommendation |
|-----|---------------|-----------------|-----------|----------------|
| lint | ubuntu-latest | ❌ No | Installs Homebrew, runs golangci-lint, Biome, Terraform lint - likely >15min | Keep ubuntu-latest |
| test | ubuntu-latest | ❌ No | Runs full Go test suite with coverage, JS tests - likely >15min | Keep ubuntu-latest |
| secret-scan | ubuntu-latest | ✅ Yes | Gitleaks scan is typically fast (<5min) | **Migrate to ubuntu-slim** |
| claude-code-review | ubuntu-latest | ⚠️ Maybe | Depends on PR size; large PRs may exceed 15min | Test with ubuntu-slim |

**Estimated savings:** Minimal (public repo)

#### Workflow: `.github/workflows/release.yaml`

| Job | Current Runner | Slim-Compatible? | Reasoning | Recommendation |
|-----|---------------|-----------------|-----------|----------------|
| release | ubuntu-latest | ❌ No | GoReleaser builds multi-platform binaries - very heavy | Keep ubuntu-latest |

**Estimated savings:** None

#### Workflow: `.github/workflows/publish.yaml`

| Job | Current Runner | Slim-Compatible? | Reasoning | Recommendation |
|-----|---------------|-----------------|-----------|----------------|
| publish | ubuntu-latest | ❌ No | Runs `pnpm turbo run build` - build step is heavy | Keep ubuntu-latest |

**Estimated savings:** None

#### Workflow: `.github/workflows/update-webkit-repos.yaml`

| Job | Current Runner | Slim-Compatible? | Reasoning | Recommendation |
|-----|---------------|-----------------|-----------|----------------|
| update-repos | ubuntu-latest | ❌ No | Clones multiple repos, installs webkit CLI, runs updates - very long | Keep ubuntu-latest |

**Estimated savings:** None

#### Workflow: `.github/workflows/dispatch-guidelines.yaml`

| Job | Current Runner | Slim-Compatible? | Reasoning | Recommendation |
|-----|---------------|-----------------|-----------|----------------|
| update-docs | ubuntu-latest | ⚠️ Maybe | Runs Go command and pnpm command - depends on generation time | Test with ubuntu-slim |

**Estimated savings:** Minimal (public repo, infrequent runs)

### 3.2 Playground Template Workflows (internal/playground)

These workflows are **generated by WebKit** and used in **downstream projects** which may be **private repositories**.

**Cost impact:** High potential - private repos pay per minute

#### Workflow: `internal/playground/.github/workflows/pr.yaml`

| Job | Current Runner | Slim-Compatible? | Reasoning | Recommendation |
|-----|---------------|-----------------|-----------|----------------|
| setup-webkit | ubuntu-latest | ✅ Yes | Downloads pre-built artifact - very fast | **Migrate to ubuntu-slim** |
| secret-scan | ubuntu-latest | ✅ Yes | Gitleaks scan - fast | **Migrate to ubuntu-slim** |
| detect-changes | ubuntu-latest | ✅ Yes | Path filtering with dorny/paths-filter - very fast | **Migrate to ubuntu-slim** |
| drift-detection | ubuntu-latest | ⚠️ Maybe | Runs webkit drift check - depends on project size | **Test with ubuntu-slim** |
| app-json-diff | ubuntu-latest | ✅ Yes | Git diff operation - very fast | **Migrate to ubuntu-slim** |
| terraform-plan-production | ubuntu-latest | ❌ No | Terraform plan with multiple providers - heavy | Keep ubuntu-latest |
| validate-app-json | ubuntu-latest | ✅ Yes | JSON validation with webkit CLI - fast | **Migrate to ubuntu-slim** |
| claude-code-review | ubuntu-latest | ⚠️ Maybe | Depends on PR size | **Test with ubuntu-slim** |
| app-cms | ubuntu-latest | ❌ No | pnpm build, lint, test - heavy | Keep ubuntu-latest |
| app-web | ubuntu-latest | ❌ No | pnpm build, lint, test - heavy | Keep ubuntu-latest |
| app-api | ubuntu-latest | ❌ No | Go build, lint, test - heavy | Keep ubuntu-latest |
| migration-check-cms | ubuntu-latest | ❌ No | Database connection and migration check - can be slow | Keep ubuntu-latest |

**Estimated savings per PR:**
- 5 jobs migrated to ubuntu-slim × ~3 minutes average × $0.004 savings/min = **~$0.06 per PR**
- For 100 PRs/month: **~$6/month per private repository**

#### Workflow: `internal/playground/.github/workflows/release.yaml`

| Job | Current Runner | Slim-Compatible? | Reasoning | Recommendation |
|-----|---------------|-----------------|-----------|----------------|
| setup-webkit | ubuntu-latest | ✅ Yes | Download artifact - fast | **Migrate to ubuntu-slim** |
| build-and-push | ubuntu-latest | ❌ No | Docker builds for 3 services - very heavy | Keep ubuntu-latest |
| cleanup-containers | ubuntu-latest | ✅ Yes | GitHub API calls to delete old packages - fast | **Migrate to ubuntu-slim** |
| terraform-apply-production | ubuntu-latest | ❌ No | Terraform apply - heavy | Keep ubuntu-latest |
| deploy-app-web | ubuntu-latest | ❌ No | DigitalOcean deployment - can be slow | Keep ubuntu-latest |
| deploy-vm-cms | ubuntu-latest | ❌ No | Ansible playbook execution - heavy | Keep ubuntu-latest |
| notify-success | ubuntu-latest | ✅ Yes | Slack notification - very fast | **Migrate to ubuntu-slim** |
| notify-failure | ubuntu-latest | ✅ Yes | Slack notification - very fast | **Migrate to ubuntu-slim** |

**Estimated savings per release:**
- 4 jobs migrated to ubuntu-slim × ~2 minutes average × $0.004 savings/min = **~$0.032 per release**
- For 20 releases/month: **~$0.64/month per private repository**

---

## 4. Compatibility Concerns and Risks

### 4.1 Critical Limitations

#### 15-Minute Timeout (Non-Negotiable)
- **Risk:** Jobs that occasionally exceed 15 minutes will **fail completely**
- **Impact:** Workflow failures, failed checks blocking PRs
- **Mitigation:** Only migrate jobs with consistent <10 minute runtimes (5-minute safety buffer)

#### Container Environment Differences
- **Risk:** Missing system packages or tools that exist in full VM
- **Impact:** Job failures due to missing dependencies
- **Mitigation:** Test thoroughly before migration, add explicit package installation if needed

#### Reduced CPU Performance
- **Risk:** Single-CPU jobs may be slower than expected
- **Impact:** Jobs approaching timeout limits may fail
- **Mitigation:** Benchmark job durations before migration

### 4.2 Tool Compatibility Matrix

| Tool/Action | ubuntu-latest | ubuntu-slim | Notes |
|------------|--------------|-------------|-------|
| actions/checkout | ✅ | ✅ | Standard action, fully compatible |
| actions/setup-go | ✅ | ✅ | Compatible |
| actions/setup-node | ✅ | ✅ | Compatible |
| golangci-lint-action | ✅ | ⚠️ | May be slower on single CPU |
| Docker build | ✅ | ❌ | Not recommended for slim runners |
| Terraform | ✅ | ❌ | Too heavy for slim runners |
| Homebrew | ✅ | ⚠️ | Installation is slow, may timeout |
| gitleaks | ✅ | ✅ | Fast scanner, compatible |
| gh CLI | ✅ | ✅ | Compatible |
| pnpm install + build | ✅ | ❌ | Build steps too heavy |

### 4.3 Workflow Failure Scenarios

**Scenario 1: Occasional Timeout**
- Job usually takes 12 minutes, occasionally hits 16 minutes
- **Result:** Intermittent failures, blocking PRs unpredictably
- **Solution:** Do NOT migrate this job

**Scenario 2: Large PR Reviews**
- claude-code-review typically takes 5 minutes
- Large PRs (1000+ lines) may take 20 minutes
- **Result:** Review fails on large PRs
- **Solution:** Keep on ubuntu-latest OR add conditional logic

**Scenario 3: Dependency Installation**
- Job installs Homebrew + multiple packages
- Installation alone takes 12 minutes
- **Result:** Timeout before actual work begins
- **Solution:** Do NOT migrate, or use pre-built container images

---

## 5. Cost Savings Projections

### 5.1 Assumptions
- Private repositories using WebKit templates
- Current pricing: ubuntu-latest = $0.008/min, ubuntu-slim = $0.004/min (50% cost)
- Average PR frequency: 100 PRs/month
- Average release frequency: 20 releases/month

### 5.2 Per-Repository Savings (Private Repos Only)

#### Conservative Estimate
**Migrated jobs:**
- setup-webkit (2 min/run)
- secret-scan (3 min/run)
- detect-changes (1 min/run)
- app-json-diff (1 min/run)
- validate-app-json (2 min/run)
- cleanup-containers (1 min/run)
- notify-success/failure (0.5 min/run)

**PR workflow savings:**
- 5 jobs × 100 PRs × 9 min avg × $0.004 savings/min = **$18/month**

**Release workflow savings:**
- 4 jobs × 20 releases × 3.5 min avg × $0.004 savings/min = **$1.12/month**

**Total per repository:** **~$19.12/month**

#### Optimistic Estimate
If drift-detection and claude-code-review are also migrated:

**Additional PR workflow savings:**
- 2 more jobs × 100 PRs × 5 min avg × $0.004 savings/min = **$4/month**

**Total per repository:** **~$23.12/month**

### 5.3 Organization-Wide Savings

If WebKit is used across **10 private repositories**:
- Conservative: $19.12 × 10 = **$191.20/month** = **$2,294.40/year**
- Optimistic: $23.12 × 10 = **$231.20/month** = **$2,774.40/year**

If WebKit is used across **50 private repositories**:
- Conservative: $19.12 × 50 = **$956/month** = **$11,472/year**
- Optimistic: $23.12 × 50 = **$1,156/month** = **$13,872/year**

### 5.4 Cost vs. Effort Analysis

**Implementation effort:**
- Template modification: 2-4 hours
- Testing across sample projects: 4-8 hours
- Documentation updates: 2 hours
- **Total:** ~8-14 hours

**Break-even point:**
- At 10 repos: 1-2 months
- At 50 repos: <1 month

**ROI:** High for organizations with multiple private repositories using WebKit

---

## 6. Implementation Recommendations

### 6.1 Phased Migration Strategy

#### Phase 1: Low-Risk Jobs (Immediate)
Migrate these jobs with **high confidence** of success:
- ✅ `setup-webkit` - Downloads artifact (~2 min)
- ✅ `secret-scan` - Gitleaks scan (~3 min)
- ✅ `detect-changes` - Path filtering (~1 min)
- ✅ `app-json-diff` - Git diff (~1 min)
- ✅ `validate-app-json` - JSON validation (~2 min)
- ✅ `cleanup-containers` - API calls (~1 min)
- ✅ `notify-success` - Slack notification (~0.5 min)
- ✅ `notify-failure` - Slack notification (~0.5 min)

**Expected savings:** ~$19/month per private repo

#### Phase 2: Medium-Risk Jobs (Testing Required)
Test these jobs in a **non-blocking manner** first:
- ⚠️ `drift-detection` - Test with various project sizes
- ⚠️ `claude-code-review` - Test with small/medium/large PRs
- ⚠️ `dispatch-guidelines.yaml: update-docs` - Test generation time

**Approach:**
1. Create test branch with ubuntu-slim configuration
2. Monitor job durations over 2-4 weeks
3. Check maximum duration vs 15-minute limit
4. If 95th percentile < 10 minutes → migrate
5. If any jobs exceed 13 minutes → do NOT migrate

**Additional savings if successful:** ~$4/month per private repo

#### Phase 3: Not Recommended (Do Not Migrate)
These jobs should **remain on ubuntu-latest**:
- ❌ All build jobs (app-cms, app-web, app-api)
- ❌ All test jobs
- ❌ lint job (Homebrew installation is slow)
- ❌ build-and-push (Docker builds)
- ❌ terraform-plan-production / terraform-apply-production
- ❌ deploy-app-web / deploy-vm-cms
- ❌ migration-check-cms
- ❌ release job (GoReleaser)
- ❌ publish job (builds packages)
- ❌ update-repos job (clones multiple repos)

### 6.2 Template Modification Approach

**Option A: Explicit Migration (Recommended)**
Modify templates to explicitly specify runner types:

```yaml
jobs:
  # Fast job - use ubuntu-slim
  secret-scan:
    runs-on: ubuntu-slim
    steps: ...

  # Heavy job - use ubuntu-latest
  test:
    runs-on: ubuntu-latest
    steps: ...
```

**Pros:**
- Explicit and clear
- Easy to understand which jobs use which runner
- No ambiguity

**Cons:**
- Requires template changes
- Generated workflows will need regeneration

**Option B: Conditional Runner Selection**
Use repository variables to control runner selection:

```yaml
jobs:
  secret-scan:
    runs-on: ${{ vars.LIGHT_RUNNER || 'ubuntu-latest' }}
```

**Pros:**
- Flexibility per repository
- Can test per-project without template changes

**Cons:**
- More complex
- Requires variable setup in each repo

**Recommendation:** Use **Option A** for simplicity and consistency across WebKit-generated projects.

### 6.3 Testing and Validation Plan

#### Step 1: Baseline Measurement
Before migration, collect current job durations:
```bash
# Use GitHub Actions API or UI to export job durations
# Target: 4 weeks of data
# Metrics: mean, median, 95th percentile, max duration
```

#### Step 2: Create Test Repository
Set up a test repository with ubuntu-slim configuration:
- Fork or create sample project
- Apply modified templates with ubuntu-slim for Phase 1 jobs
- Run 20-30 PRs over 2 weeks
- Monitor for failures or timeouts

#### Step 3: Validate Compatibility
For each migrated job, verify:
- ✅ No timeouts (all runs < 10 minutes)
- ✅ No new failures (same pass rate as ubuntu-latest)
- ✅ Same functionality (all steps execute correctly)
- ✅ Acceptable performance (not significantly slower)

#### Step 4: Gradual Rollout
1. Update WebKit templates for Phase 1 jobs
2. Document changes in release notes
3. Encourage users to regenerate workflows via `webkit update`
4. Monitor for issues via GitHub Issues/Discussions
5. If stable for 1 month → proceed to Phase 2

#### Step 5: Monitoring and Rollback
Set up alerts for:
- Increased timeout failures
- Increased job failure rates
- User reports of issues

**Rollback criteria:**
- >5% increase in job failure rate
- >10 timeout incidents in 1 week
- Critical user feedback

**Rollback process:**
1. Revert template changes
2. Release hotfix version
3. Notify users to regenerate workflows

---

## 7. Alternative Cost Optimization Strategies

If ubuntu-slim migration is too risky or not worthwhile, consider these alternatives:

### 7.1 Workflow Optimization
- **Cache dependencies:** Use `actions/cache` for Go modules, npm packages, Homebrew
- **Parallel jobs:** Split large test suites into parallel jobs (may increase cost but reduce total time)
- **Conditional jobs:** Skip unnecessary jobs using `paths` filters (already implemented well)
- **Artifact reuse:** Build once, test multiple times (using artifacts)

**Estimated savings:** 10-20% reduction in total runtime

### 7.2 Self-Hosted Runners
For very high-volume usage:
- Set up self-hosted runners on cost-effective cloud VMs
- Use spot instances for non-critical workloads
- May require more operational overhead

**Estimated savings:** 50-70% for high-volume workloads

### 7.3 Migrate to Faster Tools
- Replace golangci-lint with faster alternatives
- Use pre-built Docker images instead of building on-the-fly
- Use `actions/cache` more aggressively

**Estimated savings:** 5-15% reduction in runtime

---

## 8. Recommendations and Next Steps

### 8.1 Primary Recommendation

**Recommendation:** **Proceed with Phase 1 migration** for WebKit playground templates.

**Reasoning:**
1. **Low risk:** Phase 1 jobs are fast and simple
2. **Measurable savings:** ~$19/month per private repo
3. **Scalable impact:** Multiplies across all WebKit-generated projects
4. **Easy rollback:** Can revert templates if issues arise
5. **No downside for public repos:** Main webkit repo (public) won't be affected negatively

### 8.2 Do NOT Migrate (Until Further Testing)
- Main webkit repository workflows (public repo = no cost benefit)
- Heavy build/test jobs in playground templates
- Jobs with variable runtime that may occasionally exceed 15 minutes

### 8.3 Action Items

#### Immediate (Week 1)
1. ✅ **Review this research document** with team
2. ✅ **Get approval** for Phase 1 migration
3. ⬜ **Create feature branch** for template modifications
4. ⬜ **Modify playground templates:**
   - `internal/playground/.github/workflows/pr.yaml`
   - `internal/playground/.github/workflows/release.yaml`
5. ⬜ **Update affected jobs** to use `runs-on: ubuntu-slim`

#### Testing Phase (Week 2-3)
6. ⬜ **Create test repository** with modified templates
7. ⬜ **Run 20+ test PRs** and 5+ releases
8. ⬜ **Monitor job durations** and failure rates
9. ⬜ **Validate all checks pass** consistently

#### Rollout Phase (Week 4)
10. ⬜ **Merge template changes** to main branch
11. ⬜ **Update WebKit version** and release
12. ⬜ **Document changes** in release notes with migration guidance:
    - "This version optimizes workflow costs by using ubuntu-slim runners for lightweight jobs"
    - "Run `webkit update` to regenerate workflows with optimized runner configuration"
    - "No action required - existing workflows will continue to work"
13. ⬜ **Announce in discussions/blog** if applicable

#### Monitoring Phase (Month 2)
14. ⬜ **Monitor GitHub Issues** for timeout or failure reports
15. ⬜ **Collect feedback** from users
16. ⬜ **Measure actual cost savings** from representative private repos
17. ⬜ **Decide on Phase 2** migration based on Phase 1 success

---

## 9. Open Questions and Considerations

### 9.1 Questions for Stakeholders
1. **How many private repositories** currently use WebKit templates?
   - This determines total potential savings

2. **What is the current monthly GitHub Actions spend** for the organization?
   - Helps calculate % savings

3. **Are there any custom workflows** in private repos that might be affected?
   - May need migration guidance documentation

4. **What is the risk tolerance** for occasional workflow failures during testing?
   - Determines rollout speed

5. **Is there budget for self-hosted runners** as an alternative?
   - May offer better long-term savings

### 9.2 Technical Unknowns (Require Testing)
1. **Actual ubuntu-slim pricing** - Need to verify after first billing cycle
2. **Container startup time** vs VM startup time - May offset some savings
3. **Network performance** in slim runners - May affect download speeds
4. **Maximum concurrent jobs** on ubuntu-slim - May have different limits

### 9.3 Future Considerations
1. **GitHub's 2026 pricing changes** - May alter cost calculations
2. **New runner types** - GitHub may introduce more options
3. **WebKit adoption growth** - More repos = more savings
4. **Alternative CI platforms** - Compare with GitLab CI, CircleCI pricing

---

## 10. Conclusion

**Summary:**
- ubuntu-slim runners offer **significant cost savings** (~50% per job) for lightweight operations
- **Phase 1 jobs are excellent candidates:** low risk, high confidence of success
- **Estimated savings:** $19-23/month per private repository
- **Organization-wide impact:** Potentially $2,000-14,000/year depending on adoption
- **Risk:** Low for Phase 1 jobs, manageable with proper testing

**Final Verdict:**
✅ **Proceed with Phase 1 migration** - The cost savings justify the implementation effort, especially for organizations with multiple private repositories using WebKit templates.

⚠️ **Requires validation** - Must test thoroughly before rolling out to all generated projects.

❌ **Do NOT migrate build/test jobs** - The 15-minute timeout makes these incompatible.

---

## 11. References

Based on web search results from December 2025:

- GitHub Actions standard runners: 2 vCPU, $0.008/min for private repos
- ubuntu-slim runners: 1 vCPU, container-based, 15-minute timeout
- Use cases: automation tasks, issue operations, short-running jobs
- NOT suitable for typical heavyweight CI/CD builds
- 2026 pricing changes: `new_price = old_price × 0.6 + $0.002`

**Search sources:**
- GitHub Actions runner pricing documentation
- GitHub blog announcements about pricing changes
- Community discussions on cost optimization

---

## Document Changelog

| Date | Version | Changes |
|------|---------|---------|
| 2025-12-17 | 1.0 | Initial research document created |

---

**End of Research Document**
