<script>
  import { onMount } from 'svelte';
  import { SaveRiskAssessment, GetLatestRiskAssessment } from '../../wailsjs/go/main/App.js';

  // Risk assessment state
  let assessment = {
    id: 0,
    date: new Date().toISOString().split('T')[0],
    emotional: 0,
    fomo: 0,
    bias: 0,
    physical: 0,
    pnl: 0,
    overallScore: 0
  };

  let saving = false;
  let message = '';
  let messageType = 'info';

  // Calculate the visual position size based on risk score
  $: positionSizePercentage = calculatePositionSize(assessment.overallScore);
  $: positionSizeColor = getPositionSizeColor(assessment.overallScore);

  onMount(async () => {
    try {
      const latestAssessment = await GetLatestRiskAssessment();
      if (latestAssessment) {
        assessment = latestAssessment;
      }
    } catch (err) {
      console.error('Error fetching latest assessment:', err);
      // It's okay if there's no assessment yet
    }
  });

  // Save the current assessment
  async function saveAssessment() {
    saving = true;
    message = '';
    
    try {
      const result = await SaveRiskAssessment(assessment);
      assessment = result;
      message = 'Risk assessment saved successfully';
      messageType = 'success';
    } catch (err) {
      console.error('Error saving assessment:', err);
      message = 'Failed to save assessment';
      messageType = 'error';
    } finally {
      saving = false;
      
      // Clear message after 3 seconds
      setTimeout(() => {
        message = '';
      }, 3000);
    }
  }

  // Calculate the recommended position size as a percentage
  function calculatePositionSize(score) {
    // Simple algorithm: score of -3 = 0%, score of +3 = 100%, linear in between
    // We limit from 10% to 100% for visual purposes
    const basePercentage = ((score + 3) / 6) * 100;
    return Math.max(10, Math.min(100, basePercentage));
  }

  // Get a color for the position size indicator
  function getPositionSizeColor(score) {
    if (score < -1) return '#e74c3c'; // Red for negative scores
    if (score < 1) return '#f39c12';  // Yellow for neutral scores
    return '#2ecc71';                 // Green for positive scores
  }
</script>

<div class="risk-dashboard">
  <h2>Risk Management Dashboard</h2>
  <p class="dashboard-description">
    Assess your daily emotional and psychological state to determine optimal position sizing.
  </p>

  {#if message}
    <div class="message {messageType}">
      {message}
    </div>
  {/if}

  <div class="assessment-form">
    <div class="date-selection">
      <label for="assessment-date">Assessment Date:</label>
      <input 
        type="date" 
        id="assessment-date"
        bind:value={assessment.date}
      />
    </div>

    <div class="sliders">
      <div class="slider-group">
        <label for="emotional">Emotional State: {assessment.emotional}</label>
        <input 
          type="range" 
          id="emotional" 
          min="-3" 
          max="3" 
          step="1" 
          bind:value={assessment.emotional}
        />
        <div class="slider-labels">
          <span>Negative (-3)</span>
          <span>Neutral (0)</span>
          <span>Positive (+3)</span>
        </div>
        <div class="slider-description">
          How are you feeling emotionally today? Negative values indicate stress, anxiety, or emotional disturbance.
        </div>
      </div>

      <div class="slider-group">
        <label for="fomo">FOMO Level: {assessment.fomo}</label>
        <input 
          type="range" 
          id="fomo" 
          min="-3" 
          max="3" 
          step="1" 
          bind:value={assessment.fomo}
        />
        <div class="slider-labels">
          <span>Low (-3)</span>
          <span>Medium (0)</span>
          <span>High (+3)</span>
        </div>
        <div class="slider-description">
          How strongly are you feeling the "Fear Of Missing Out"? Higher values indicate stronger FOMO which can lead to impulsive decisions.
        </div>
      </div>

      <div class="slider-group">
        <label for="bias">Market Bias: {assessment.bias}</label>
        <input 
          type="range" 
          id="bias" 
          min="-3" 
          max="3" 
          step="1" 
          bind:value={assessment.bias}
        />
        <div class="slider-labels">
          <span>Bearish (-3)</span>
          <span>Neutral (0)</span>
          <span>Bullish (+3)</span>
        </div>
        <div class="slider-description">
          Are you biased toward a bearish or bullish perspective? Strong biases can affect your trading decisions.
        </div>
      </div>

      <div class="slider-group">
        <label for="physical">Physical Condition: {assessment.physical}</label>
        <input 
          type="range" 
          id="physical" 
          min="-3" 
          max="3" 
          step="1" 
          bind:value={assessment.physical}
        />
        <div class="slider-labels">
          <span>Poor (-3)</span>
          <span>Average (0)</span>
          <span>Excellent (+3)</span>
        </div>
        <div class="slider-description">
          How is your physical well-being today? Fatigue, illness, or poor sleep can impair decision-making.
        </div>
      </div>

      <div class="slider-group">
        <label for="pnl">Recent P&L Impact: {assessment.pnl}</label>
        <input 
          type="range" 
          id="pnl" 
          min="-3" 
          max="3" 
          step="1" 
          bind:value={assessment.pnl}
        />
        <div class="slider-labels">
          <span>Losses (-3)</span>
          <span>Breakeven (0)</span>
          <span>Profits (+3)</span>
        </div>
        <div class="slider-description">
          How have your recent trading results affected your mindset? Recent losses can lead to revenge trading.
        </div>
      </div>
    </div>

    <div class="actions">
      <button class="save-button" on:click={saveAssessment} disabled={saving}>
        {saving ? 'Saving...' : 'Save Assessment'}
      </button>
    </div>

    <div class="results-panel">
      <h3>Recommended Position Size</h3>
      <div class="position-size-container">
        <div 
          class="position-size-indicator" 
          style="width: {positionSizePercentage}%; background-color: {positionSizeColor};"
        >
          {Math.round(positionSizePercentage)}%
        </div>
      </div>
      <div class="position-advice">
        {#if positionSizePercentage < 30}
          <p><strong>Warning:</strong> Consider smaller position sizes or no trading today.</p>
          <ul>
            <li>Your current state suggests elevated trading risk</li>
            <li>Consider paper trading or taking a day off</li>
            <li>If you must trade, reduce position size by 70-80%</li>
          </ul>
        {:else if positionSizePercentage < 70}
          <p><strong>Caution:</strong> Standard position sizing recommended with care.</p>
          <ul>
            <li>Trade with additional awareness of your current limitations</li>
            <li>Consider reducing position size by 30-50%</li>
            <li>Focus on high-conviction setups only</li>
          </ul>
        {:else}
          <p><strong>Optimal:</strong> Conditions favorable for your standard position size.</p>
          <ul>
            <li>Your mental and physical state supports good decision making</li>
            <li>Follow your trading plan with confidence</li>
            <li>Still maintain discipline and risk management</li>
          </ul>
        {/if}
      </div>
    </div>
  </div>

  <div class="risk-guidelines">
    <h3>Daily Trading Psychology Guidelines</h3>
    <div class="guidelines-content">
      <div class="guideline">
        <h4>✓ Start with Self-Assessment</h4>
        <p>Begin each trading day by honestly evaluating your mental and physical state.</p>
      </div>
      <div class="guideline">
        <h4>✓ Adjust Position Sizing</h4>
        <p>Use your risk assessment score to modify position sizing - trade smaller on difficult days.</p>
      </div>
      <div class="guideline">
        <h4>✓ Recognize Emotional Triggers</h4>
        <p>Be aware of market events or conditions that might trigger emotional responses.</p>
      </div>
      <div class="guideline">
        <h4>✓ Implement Circuit Breakers</h4>
        <p>Have predefined rules for when to stop trading (e.g., after 2-3 consecutive losses).</p>
      </div>
    </div>
  </div>
</div>

<style>
  .risk-dashboard {
    max-width: 800px;
    margin: 0 auto;
  }

  h2 {
    color: var(--text-primary);
    border-bottom: 2px solid var(--accent-color, #3498db);
    padding-bottom: 0.5rem;
  }

  h3 {
    color: var(--text-primary);
    margin-top: 0;
    margin-bottom: 1rem;
  }

  h4 {
    color: var(--text-primary);
    margin: 0 0 0.5rem 0;
  }

  .dashboard-description {
    color: var(--text-secondary);
    margin-bottom: 2rem;
  }

  .assessment-form {
    background-color: var(--bg-primary);
    border-radius: 8px;
    box-shadow: 0 2px 10px var(--shadow-color);
    padding: 2rem;
    margin-bottom: 2rem;
  }

  .date-selection {
    margin-bottom: 1.5rem;
  }

  .date-selection input {
    padding: 0.5rem;
    border: 1px solid var(--border-color);
    border-radius: 4px;
    background-color: var(--bg-primary);
    color: var(--text-primary);
  }

  .sliders {
    display: grid;
    gap: 1.5rem;
  }

  .slider-group {
    margin-bottom: 1rem;
  }

  .slider-group label {
    display: block;
    margin-bottom: 0.5rem;
    font-weight: bold;
    color: var(--text-primary);
  }

  .slider-group input[type="range"] {
    width: 100%;
    margin: 0.5rem 0;
  }

  .slider-labels {
    display: flex;
    justify-content: space-between;
    font-size: 0.8rem;
    color: var(--text-secondary);
  }

  .slider-description {
    font-size: 0.85rem;
    color: var(--text-secondary);
    margin-top: 0.5rem;
    font-style: italic;
  }

  .actions {
    margin-top: 2rem;
    text-align: center;
  }

  .save-button {
    background-color: var(--accent-color, #3498db);
    color: white;
    border: none;
    padding: 0.8rem 2rem;
    border-radius: 4px;
    font-size: 1rem;
    cursor: pointer;
    transition: background-color 0.3s;
  }

  .save-button:hover {
    background-color: var(--accent-hover, #2980b9);
  }

  .save-button:disabled {
    background-color: var(--neutral-color, #95a5a6);
    cursor: not-allowed;
  }

  .message {
    padding: 1rem;
    margin-bottom: 1rem;
    border-radius: 4px;
    text-align: center;
  }

  .message.success {
    background-color: var(--success-bg, #d5f5e3);
    color: var(--success-color, #27ae60);
  }

  .message.error {
    background-color: var(--error-bg, #fadbd8);
    color: var(--error-color, #e74c3c);
  }

  .message.info {
    background-color: var(--info-bg, #d6eaf8);
    color: var(--accent-color, #3498db);
  }

  .results-panel {
    margin-top: 2rem;
    padding-top: 1.5rem;
    border-top: 1px solid var(--border-color);
  }

  .position-size-container {
    background-color: var(--bg-secondary);
    border-radius: 4px;
    height: 30px;
    overflow: hidden;
    margin: 1rem 0;
  }

  .position-size-indicator {
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    color: white;
    font-weight: bold;
    transition: width 0.5s, background-color 0.5s;
  }

  .position-advice {
    margin-top: 1rem;
    padding: 1rem;
    border-radius: 8px;
    background-color: var(--bg-secondary);
  }

  .position-advice p {
    margin-top: 0;
  }

  .position-advice ul {
    margin-bottom: 0;
  }

  .risk-guidelines {
    background-color: var(--bg-primary);
    border-radius: 8px;
    box-shadow: 0 2px 10px var(--shadow-color);
    padding: 2rem;
  }

  .guidelines-content {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 1.5rem;
  }

  .guideline {
    padding: 1rem;
    border-radius: 8px;
    background-color: var(--bg-secondary);
  }

  .guideline h4 {
    color: var(--success-color, #27ae60);
  }

  .guideline p {
    margin: 0;
    font-size: 0.9rem;
    color: var(--text-primary);
  }

  @media (max-width: 768px) {
    .guidelines-content {
      grid-template-columns: 1fr;
    }
  }
</style> 