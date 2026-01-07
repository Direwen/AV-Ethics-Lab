# LLM Production Strategy - Don't Break the Bank or the Experiment

## TL;DR
We're using free Groq accounts with smart caching and fallbacks to handle experimental data collection without compromising research quality. No self-hosted models because they're slow AF.

## The Setup

### API Key Strategy
- **3 Groq accounts** (max we can manage right now)
- **43,200 requests/day total** (14,400 per account)
- **1,800 requests/hour peak** (600 per account)
- **Round-robin rotation** with usage tracking

### The Fallback Hierarchy (Quality First)

1. **Cache Hit (70% of requests)** - Instant win
   - Semantic similarity matching
   - Same quality as fresh LLM
   - No API calls needed

2. **Fresh LLM Generation (25% of requests)** - The good stuff
   - Queue system with 3 Groq keys
   - 3-5 second response time
   - Perfect research quality

3. **Pre-generated Pool (4% of requests)** - Our insurance policy
   - Batch-generated during off-hours using same LLM
   - Same quality, just made earlier
   - Fallback when APIs are completely dead

4. **Relaxed Cache (0.9% of requests)** - Slightly different but still valid
   - Lower similarity threshold (75% vs 95%)
   - Minor variations in position/factors
   - Still research-worthy

5. **Graceful Failure (0.1% of requests)** - Better than bad data
   - "Try again in a few minutes"
   - Maintains research integrity
   - No compromised scenarios

## Why This Works for Research

### Data Quality Priorities
- **Never** use template-based generation (too predictable)
- **Never** compromise on LLM quality for core scenarios
- **Always** prefer failure over bad data
- **Cache variations** are acceptable (same LLM origin)

### Expected Performance
- **Normal days**: 99% success rate, high quality
- **API outages**: 99% success rate, mostly high quality
- **Complete disaster**: 70% success rate (cache only)

## Implementation Strategy

### Phase 1: Core Features First
Focus on getting the basic flow working:
- Single Groq account
- Basic caching
- Simple error handling

### Phase 2: Production Hardening (Later)
When ready to scale:
- Multi-account rotation
- Request queueing system
- Pre-generated pool
- Relaxed cache matching
- Monitoring and alerts

## Technical Details (For Later)

### Request Queueing Concept
- **Non-blocking**: API returns immediately with "generating..." status
- **Polling**: Frontend checks status every 500ms
- **Session isolation**: Users can only access their own scenarios
- **Error resilience**: Circuit breaker prevents cascade failures

### Cache Strategy
- **Semantic similarity**: Hash based on template + factors + position
- **Hit rate target**: 60-70% (higher = too repetitive, lower = inefficient)
- **Variations**: Slight modifications to cached scenarios for variety

### Pre-generated Pool
- **Batch generation**: Run during 3-6 AM when quota is fresh
- **Popular combinations**: Focus on common template/factor pairs
- **Storage**: Database with metadata for matching
- **Refresh cycle**: Weekly regeneration to keep content fresh

### Monitoring Must-haves
- Cache hit rates
- API key usage per account
- Queue depth and processing time
- Error rates and types
- User satisfaction (completed vs failed requests)

## Cost Estimates

### With 70% Cache Hit Rate
- **Effective capacity**: ~144,000 scenarios/day
- **Concurrent users**: Hundreds (depending on usage patterns)
- **Cost**: $0 (using free tiers)

### Scaling Triggers
- If cache hit rate drops below 60%: Improve similarity matching
- If API limits hit regularly: Add more Groq accounts
- If queue depth > 100: Consider paid tiers or optimization

## Fallback Quality Assessment

### Acceptable for Research
- ✅ Fresh LLM generation
- ✅ Cached scenarios (same LLM origin)
- ✅ Pre-generated pool (same LLM, batch processed)
- ✅ Relaxed cache matches (minor variations)

### Not Acceptable for Research
- ❌ Template-based generation (too predictable)
- ❌ Heavily modified cached scenarios
- ❌ Cross-contaminated scenarios from other experiments

## Risk Mitigation

### What Could Go Wrong
- **All Groq accounts exhausted**: Use pre-generated pool
- **Groq API completely down**: Cache + pre-generated scenarios
- **Cache corruption**: Rebuild from database
- **Queue system failure**: Direct API calls with rate limiting
- **Database issues**: In-memory cache as backup

### Monitoring Alerts
- API key usage > 90% of daily limit
- Cache hit rate < 50%
- Queue processing time > 10 seconds
- Error rate > 5%
- Pre-generated pool < 100 scenarios

## Notes for Future Me

### Don't Forget
- Test the queue system thoroughly (edge cases are nasty)
- Monitor cache hit rates religiously
- Keep pre-generated pool fresh
- Document all fallback scenarios
- Set up proper alerting before going live

### Nice to Have Later
- A/B testing different cache strategies
- Machine learning for better similarity matching
- Predictive pre-generation based on user patterns
- Multi-region deployment for better latency

### Definitely Don't Do
- Self-hosted models (too slow for real-time)
- Complex prompt engineering (keep it simple)
- Over-optimization before we have real usage data
- Compromising research quality for uptime

---

**Bottom Line**: This strategy gives us 99%+ uptime with research-quality data while staying on free tiers. Focus on core features first, implement this when ready to scale.