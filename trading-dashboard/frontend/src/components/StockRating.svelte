<script>
  import { onMount } from 'svelte';
  import { SaveStockRating, GetAllStockRatings } from '../../wailsjs/go/main/App.js';
  import { models } from '../../wailsjs/go/models';

  // Stock rating state
  let rating = new models.StockRating({
    id: 0,
    date: new Date().toISOString().split('T')[0],
    ticker: '',
    marketSentiment: 0,
    sectorSentiment: 0,
    basicMaterials: 0,
    communicationServices: 0,
    consumerCyclical: 0,
    consumerDefensive: 0,
    energy: 0,
    financial: 0,
    healthcare: 0,
    industrials: 0,
    realEstate: 0,
    technology: 0,
    utilities: 0,
    stockSentiment: 0,
    pattern: '',
    enthusiasmRating: 0
  });

  // Track which sector is relevant to the current stock
  let selectedSector = '';

  // Sector definitions
  const sectors = [
    { id: 'basicMaterials', name: 'Basic Materials' },
    { id: 'communicationServices', name: 'Communication Services' },
    { id: 'consumerCyclical', name: 'Consumer Cyclical' },
    { id: 'consumerDefensive', name: 'Consumer Defensive' },
    { id: 'energy', name: 'Energy' },
    { id: 'financial', name: 'Financial' },
    { id: 'healthcare', name: 'Healthcare' },
    { id: 'industrials', name: 'Industrials' },
    { id: 'realEstate', name: 'Real Estate' },
    { id: 'technology', name: 'Technology' },
    { id: 'utilities', name: 'Utilities' }
  ];

  // Chart patterns
  const patterns = [
    { value: '', label: 'Select a pattern...' },
    { value: 'High Base', label: 'High Base' },
    { value: 'Low Base', label: 'Low Base' },
    { value: 'Ascending Triangle', label: 'Ascending Triangle' },
    { value: 'Descending Triangle', label: 'Descending Triangle' },
    { value: 'Bull Pullback', label: 'Bull Pullback' },
    { value: 'Bear Rally', label: 'Bear Rally' },
    { value: 'Double-Top', label: 'Double-Top' },
    { value: 'Cup-and-Handle', label: 'Cup-and-Handle' },
    // Additional patterns
    { value: 'Head and Shoulders', label: 'Head and Shoulders' },
    { value: 'Inverse Head and Shoulders', label: 'Inverse Head and Shoulders' },
    { value: 'Bullish Flag', label: 'Bullish Flag' },
    { value: 'Bearish Flag', label: 'Bearish Flag' },
    { value: 'Rising Wedge', label: 'Rising Wedge' },
    { value: 'Falling Wedge', label: 'Falling Wedge' },
    { value: 'Double Bottom', label: 'Double Bottom' },
    { value: 'Rounding Bottom', label: 'Rounding Bottom (Saucer)' },
    { value: 'Breakaway Gap', label: 'Breakaway Gap' },
    { value: 'Runaway Gap', label: 'Runaway Gap' },
    { value: 'Exhaustion Gap', label: 'Exhaustion Gap' },
    { value: 'Bullish Engulfing', label: 'Bullish Engulfing' },
    { value: 'Bearish Engulfing', label: 'Bearish Engulfing' }
  ];

  // Pattern points for visual display
  const patternPoints = {
    'High Base': 2,
    'Low Base': 2,
    'Ascending Triangle': 3,
    'Descending Triangle': 3,
    'Bull Pullback': 2,
    'Bear Rally': 2,
    'Double-Top': 3,
    'Cup-and-Handle': 4,
    // Additional patterns with point values
    'Head and Shoulders': 4,
    'Inverse Head and Shoulders': 4,
    'Bullish Flag': 3,
    'Bearish Flag': 3,
    'Rising Wedge': 2,
    'Falling Wedge': 2,
    'Double Bottom': 3,
    'Rounding Bottom': 3,
    'Breakaway Gap': 3,
    'Runaway Gap': 2,
    'Exhaustion Gap': 1,
    'Bullish Engulfing': 2,
    'Bearish Engulfing': 2
  };

  // UI state
  let saving = false;
  let message = '';
  let messageType = 'info';
  let recentRatings = [];

  // Get pattern points for the currently selected pattern
  $: currentPatternPoints = rating.pattern ? patternPoints[rating.pattern] : 0;
  
  // Calculate the enthusiasm rating client-side (just for display)
  $: enthusiasmRating = rating.stockSentiment + currentPatternPoints;

  // Get the sentiment value for the selected sector
  $: sectorSentiment = selectedSector ? rating[selectedSector] : 0;

  onMount(async () => {
    await loadRecentRatings();
  });

  // Load recent stock ratings
  async function loadRecentRatings() {
    try {
      recentRatings = await GetAllStockRatings();
    } catch (err) {
      console.error('Error fetching stock ratings:', err);
    }
  }

  // Save the current stock rating
  async function saveRating() {
    if (!rating.ticker) {
      message = 'Please enter a ticker symbol';
      messageType = 'error';
      setTimeout(() => {
        message = '';
      }, 3000);
      return;
    }

    if (!rating.pattern) {
      message = 'Please select a chart pattern';
      messageType = 'error';
      setTimeout(() => {
        message = '';
      }, 3000);
      return;
    }

    if (!selectedSector) {
      message = 'Please select a sector for this stock';
      messageType = 'error';
      setTimeout(() => {
        message = '';
      }, 3000);
      return;
    }

    saving = true;
    message = '';
    
    try {
      const result = await SaveStockRating(rating);
      rating.id = result.id;
      rating.enthusiasmRating = result.enthusiasmRating;
      
      message = `Stock rating saved successfully with enthusiasm score: ${rating.enthusiasmRating}`;
      messageType = 'success';
      
      // Reset the form for new entry but keep market and sector sentiments
      const { 
        marketSentiment, 
        basicMaterials,
        communicationServices,
        consumerCyclical,
        consumerDefensive,
        energy,
        financial,
        healthcare,
        industrials,
        realEstate,
        technology,
        utilities
      } = rating;
      
      rating = new models.StockRating({
        id: 0,
        date: new Date().toISOString().split('T')[0],
        ticker: '',
        marketSentiment,
        sectorSentiment,
        basicMaterials,
        communicationServices,
        consumerCyclical,
        consumerDefensive,
        energy,
        financial,
        healthcare,
        industrials,
        realEstate,
        technology,
        utilities,
        stockSentiment: 0,
        pattern: '',
        enthusiasmRating: 0
      });
      
      selectedSector = '';
      
      // Reload recent ratings
      await loadRecentRatings();
    } catch (err) {
      console.error('Error saving stock rating:', err);
      message = 'Failed to save stock rating';
      messageType = 'error';
    } finally {
      saving = false;
      
      // Clear message after 3 seconds
      setTimeout(() => {
        message = '';
      }, 3000);
    }
  }

  // Load an existing rating for editing
  function editRating(existingRating) {
    rating = new models.StockRating(existingRating);
    // Find the sector that has a non-zero value
    selectedSector = sectors.find(sector => existingRating[sector.id] !== 0)?.id || '';
  }

  // Reset the form
  function resetForm() {
    rating = new models.StockRating({
      id: 0,
      date: new Date().toISOString().split('T')[0],
      ticker: '',
      marketSentiment: 0,
      sectorSentiment: 0,
      basicMaterials: 0,
      communicationServices: 0,
      consumerCyclical: 0,
      consumerDefensive: 0,
      energy: 0,
      financial: 0,
      healthcare: 0,
      industrials: 0,
      realEstate: 0,
      technology: 0,
      utilities: 0,
      stockSentiment: 0,
      pattern: '',
      enthusiasmRating: 0
    });
    selectedSector = '';
  }

  // Get color based on sentiment value
  function getSentimentColor(value) {
    if (value < -1) return '#e74c3c'; // Red for negative
    if (value < 1) return '#f39c12';  // Yellow for neutral
    return '#2ecc71';                 // Green for positive
  }

  // Group patterns by category for better organization
  const patternCategories = [
    {
      name: 'Common Patterns',
      patterns: ['High Base', 'Low Base', 'Ascending Triangle', 'Descending Triangle', 
                 'Bull Pullback', 'Bear Rally', 'Double-Top', 'Cup-and-Handle']
    },
    {
      name: 'Head & Shoulders Patterns',
      patterns: ['Head and Shoulders', 'Inverse Head and Shoulders']
    },
    {
      name: 'Flag & Wedge Patterns',
      patterns: ['Bullish Flag', 'Bearish Flag', 'Rising Wedge', 'Falling Wedge']
    },
    {
      name: 'Bottom Patterns',
      patterns: ['Double Bottom', 'Rounding Bottom']
    },
    {
      name: 'Gap Patterns',
      patterns: ['Breakaway Gap', 'Runaway Gap', 'Exhaustion Gap']
    },
    {
      name: 'Engulfing Patterns',
      patterns: ['Bullish Engulfing', 'Bearish Engulfing']
    }
  ];

  // Flag to show pattern info
  let showPatternInfo = false;
  let selectedPatternForInfo = '';

  function togglePatternInfo(pattern) {
    if (selectedPatternForInfo === pattern && showPatternInfo) {
      showPatternInfo = false;
    } else {
      selectedPatternForInfo = pattern;
      showPatternInfo = true;
    }
  }

  // Get pattern description
  function getPatternDescription(pattern) {
    const descriptions = {
      'High Base': 'Consolidation pattern near resistance with tight price action, suggesting strength.',
      'Low Base': 'Consolidation pattern near support with tight price action, showing potential for reversal.',
      'Ascending Triangle': 'Bullish pattern with horizontal resistance and rising support, typically breaks upward.',
      'Descending Triangle': 'Bearish pattern with horizontal support and falling resistance, typically breaks downward.',
      'Bull Pullback': 'Temporary price retreat within an uptrend, often creating a buying opportunity.',
      'Bear Rally': 'Temporary price rise within a downtrend, potentially creating a shorting opportunity.',
      'Double-Top': 'Bearish reversal pattern showing two roughly equal highs, indicating resistance.',
      'Cup-and-Handle': 'Bullish continuation pattern resembling a cup with a handle, signaling continuation.',
      'Head and Shoulders': 'Bearish reversal pattern with three peaks (middle highest), signaling a trend change.',
      'Inverse Head and Shoulders': 'Bullish reversal pattern with three troughs (middle lowest), signaling an uptrend.',
      'Bullish Flag': 'Continuation pattern that forms after a strong upward move, followed by consolidation.',
      'Bearish Flag': 'Continuation pattern that forms after a strong downward move, followed by consolidation.',
      'Rising Wedge': 'Pattern with converging trend lines sloping upward, often breaks downward.',
      'Falling Wedge': 'Pattern with converging trend lines sloping downward, often breaks upward.',
      'Double Bottom': 'Bullish reversal pattern showing two roughly equal lows, indicating support.',
      'Rounding Bottom': 'Long-term reversal pattern indicating gradual shift from bearish to bullish sentiment.',
      'Breakaway Gap': 'Gap that forms at the beginning of a trend, signaling a strong move.',
      'Runaway Gap': 'Gap that forms during the middle of a trend, confirming the trend strength.',
      'Exhaustion Gap': 'Gap that forms near the end of a trend, signaling potential reversal.',
      'Bullish Engulfing': 'Two-candle reversal pattern where a bullish candle completely engulfs the previous bearish one.',
      'Bearish Engulfing': 'Two-candle reversal pattern where a bearish candle completely engulfs the previous bullish one.'
    };
    
    return descriptions[pattern] || 'No description available';
  }

  // Function to get the sector name by id
  function getSectorName(sectorId) {
    const sector = sectors.find(s => s.id === sectorId);
    return sector ? sector.name : '';
  }
  
  // Function to find the highest rated sector for a stock
  function getHighestRatedSector(item) {
    let highestSector = { id: '', value: -4 }; // Start with value below possible range
    
    sectors.forEach(sector => {
      if (item[sector.id] > highestSector.value) {
        highestSector = { id: sector.id, value: item[sector.id] };
      }
    });
    
    return highestSector;
  }
</script>

<div class="stock-rating">
  <h2>Stock Rating Dashboard</h2>
  
  <p class="dashboard-description">
    Rate market conditions, sectors, and individual stocks to identify the best trading opportunities.
  </p>

  {#if message}
    <div class="message {messageType}">
      {message}
    </div>
  {/if}

  <div class="rating-form">
    <div class="form-container">
      <div class="column">
        <h3>Market & Sectors</h3>
        
        <div class="date-selection">
          <label for="rating-date">Rating Date:</label>
          <input 
            type="date" 
            id="rating-date"
            bind:value={rating.date}
          />
        </div>

        <div class="slider-group">
          <label for="market-sentiment">Overall Market: {rating.marketSentiment}</label>
          <input 
            type="range" 
            id="market-sentiment" 
            min="-3" 
            max="3" 
            step="1" 
            bind:value={rating.marketSentiment}
          />
          <div class="slider-labels">
            <span>Bearish (-3)</span>
            <span>Bullish (+3)</span>
          </div>
        </div>

        <div class="sector-container">
          <div class="sector-header">
            <h4>Sector Ratings</h4>
            <div class="input-group">
              <label for="selected-sector">Primary Sector for This Stock:</label>
              <select id="selected-sector" bind:value={selectedSector}>
                <option value="">Select primary sector...</option>
                {#each sectors as sector}
                  <option value={sector.id}>{sector.name}</option>
                {/each}
              </select>
            </div>
          </div>

          <!-- Always show all sector sliders -->
          {#each sectors as sector}
            <div class="slider-group" class:highlighted-sector={selectedSector === sector.id}>
              <label for="{sector.id}">{sector.name}: {rating[sector.id]}</label>
              <input 
                type="range" 
                id="{sector.id}" 
                min="-3" 
                max="3" 
                step="1" 
                bind:value={rating[sector.id]}
              />
              <div class="slider-labels">
                <span>Bearish (-3)</span>
                <span>Bullish (+3)</span>
              </div>
            </div>
          {/each}
        </div>
      </div>

      <div class="column">
        <h3>Stock Rating</h3>

        <div class="input-group">
          <label for="ticker">Ticker Symbol:</label>
          <input 
            type="text" 
            id="ticker" 
            bind:value={rating.ticker}
            placeholder="e.g., AAPL"
            maxlength="10"
          />
        </div>

        <div class="slider-group">
          <label for="stock-sentiment">Stock Sentiment: {rating.stockSentiment}</label>
          <input 
            type="range" 
            id="stock-sentiment" 
            min="-3" 
            max="3" 
            step="1" 
            bind:value={rating.stockSentiment}
          />
          <div class="slider-labels">
            <span>Bearish (-3)</span>
            <span>Bullish (+3)</span>
          </div>
        </div>

        <div class="input-group">
          <label for="pattern">Chart Pattern:</label>
          <select id="pattern" bind:value={rating.pattern}>
            <option value="">Select a pattern...</option>
            {#each patternCategories as category}
              <optgroup label={category.name}>
                {#each category.patterns as patternValue}
                  <option value={patternValue}>{patternValue}</option>
                {/each}
              </optgroup>
            {/each}
          </select>
          <div class="pattern-help">
            <button class="info-button" on:click={() => togglePatternInfo(rating.pattern)} disabled={!rating.pattern}>
              Pattern Info
            </button>
          </div>
        </div>

        {#if showPatternInfo && selectedPatternForInfo}
          <div class="pattern-info-box">
            <h4>{selectedPatternForInfo}</h4>
            <p>{getPatternDescription(selectedPatternForInfo)}</p>
            <p class="pattern-points">Trading value: +{patternPoints[selectedPatternForInfo]} points</p>
          </div>
        {:else if rating.pattern}
          <div class="pattern-info">
            Pattern value: +{currentPatternPoints} points
          </div>
        {/if}

        <div class="enthusiasm-preview">
          <h4>Estimated Enthusiasm: <span style="color: {getSentimentColor(enthusiasmRating)}">{enthusiasmRating}</span></h4>
        </div>
      </div>
    </div>

    <div class="actions">
      <button class="reset-button" on:click={resetForm}>Reset</button>
      <button class="save-button" on:click={saveRating} disabled={saving}>
        {saving ? 'Saving...' : 'Save Rating'}
      </button>
    </div>
  </div>

  <div class="recent-ratings">
    <h3>Recent Stock Ratings</h3>
    
    {#if recentRatings.length === 0}
      <p class="no-ratings">No stock ratings have been saved yet.</p>
    {:else}
      <div class="ratings-table-container">
        <table class="ratings-table">
          <thead>
            <tr>
              <th>Date</th>
              <th>Ticker</th>
              <th>Market</th>
              <th>Primary Sector</th>
              <th>Sector Rating</th>
              <th>Stock</th>
              <th>Pattern</th>
              <th>Enthusiasm</th>
              <th>Action</th>
            </tr>
          </thead>
          <tbody>
            {#each recentRatings as item}
              {@const highestSector = getHighestRatedSector(item)}
              <tr>
                <td>{new Date(item.date).toLocaleDateString()}</td>
                <td>{item.ticker}</td>
                <td style="color: {getSentimentColor(item.marketSentiment)}">{item.marketSentiment}</td>
                <td>{getSectorName(highestSector.id)}</td>
                <td style="color: {getSentimentColor(highestSector.value)}">{highestSector.value}</td>
                <td style="color: {getSentimentColor(item.stockSentiment)}">{item.stockSentiment}</td>
                <td>{item.pattern}</td>
                <td style="color: {getSentimentColor(item.enthusiasmRating)}">{item.enthusiasmRating}</td>
                <td>
                  <button class="edit-button" on:click={() => editRating(item)}>Edit</button>
                </td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>
    {/if}
  </div>
</div>

<style>
  /* Maintain only component-specific styles */
  .stock-rating {
    max-width: 1000px;
    margin: 0 auto;
  }

  h2 {
    color: var(--text-primary);
    border-bottom: 2px solid var(--accent-color);
    padding-bottom: 0.5rem;
  }

  h3 {
    color: var(--text-primary);
    margin-top: 0;
  }

  h4 {
    margin: 0 0 0.5rem 0;
    color: var(--text-primary);
  }

  .dashboard-description {
    color: var(--text-secondary);
    margin-bottom: 2rem;
  }

  .rating-form {
    background-color: var(--bg-primary);
    border-radius: 8px;
    box-shadow: 0 2px 10px var(--shadow-color);
    padding: 2rem;
    margin-bottom: 2rem;
  }

  .form-container {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 2rem;
  }

  .column {
    display: flex;
    flex-direction: column;
  }

  .date-selection {
    margin-bottom: 1.5rem;
  }

  .date-selection input {
    padding: 0.5rem;
    border: 1px solid var(--border-color);
    border-radius: 4px;
    width: 100%;
    background-color: var(--bg-primary);
    color: var(--text-primary);
  }

  .input-group {
    margin-bottom: 1.5rem;
  }

  .input-group label {
    display: block;
    margin-bottom: 0.5rem;
    font-weight: bold;
    color: var(--text-primary);
  }

  .input-group input, .input-group select {
    width: 100%;
    padding: 0.5rem;
    border: 1px solid var(--border-color);
    border-radius: 4px;
    background-color: var(--bg-primary);
    color: var(--text-primary);
  }

  .slider-group {
    margin-bottom: 1.5rem;
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

  .pattern-info {
    font-size: 0.9rem;
    color: var(--accent-color);
    margin-bottom: 1rem;
  }

  .pattern-help {
    margin-top: 0.5rem;
    display: flex;
    justify-content: flex-end;
  }

  .info-button {
    background-color: var(--neutral-color);
    color: white;
    border: none;
    padding: 0.3rem 0.6rem;
    border-radius: 4px;
    font-size: 0.8rem;
    cursor: pointer;
  }

  .info-button:hover {
    background-color: var(--neutral-hover);
  }

  .info-button:disabled {
    background-color: var(--border-color);
    color: var(--text-secondary);
    cursor: not-allowed;
  }

  .pattern-info-box {
    background-color: var(--bg-secondary);
    padding: 1rem;
    border-radius: 4px;
    margin-bottom: 1rem;
    border-left: 4px solid var(--accent-color);
  }

  .pattern-info-box p {
    margin: 0.5rem 0;
    font-size: 0.9rem;
    color: var(--text-primary);
  }

  .pattern-points {
    color: var(--accent-color);
    font-weight: bold;
  }

  .enthusiasm-preview {
    background-color: var(--bg-secondary);
    padding: 1rem;
    border-radius: 4px;
    text-align: center;
    margin-top: 1rem;
  }

  .actions {
    margin-top: 2rem;
    display: flex;
    justify-content: flex-end;
    gap: 1rem;
  }

  .save-button {
    background-color: var(--accent-color);
    color: white;
    border: none;
    padding: 0.8rem 2rem;
    border-radius: 4px;
    font-size: 1rem;
    cursor: pointer;
    transition: background-color 0.3s;
  }

  .save-button:hover {
    background-color: var(--accent-hover);
  }

  .save-button:disabled {
    background-color: var(--neutral-color);
    cursor: not-allowed;
  }

  .reset-button {
    background-color: var(--neutral-color);
    color: white;
    border: none;
    padding: 0.8rem 2rem;
    border-radius: 4px;
    font-size: 1rem;
    cursor: pointer;
    transition: background-color 0.3s;
  }

  .reset-button:hover {
    background-color: var(--neutral-hover);
  }

  .message {
    padding: 1rem;
    margin-bottom: 1rem;
    border-radius: 4px;
    text-align: center;
  }

  .message.success {
    background-color: var(--success-bg);
    color: var(--success-color);
  }

  .message.error {
    background-color: var(--error-bg);
    color: var(--error-color);
  }

  .message.info {
    background-color: var(--info-bg);
    color: var(--accent-color);
  }

  .recent-ratings {
    background-color: var(--bg-primary);
    border-radius: 8px;
    box-shadow: 0 2px 10px var(--shadow-color);
    padding: 2rem;
  }

  .no-ratings {
    text-align: center;
    color: var(--text-secondary);
    font-style: italic;
  }

  .ratings-table-container {
    overflow-x: auto;
  }

  .ratings-table {
    width: 100%;
    border-collapse: collapse;
  }

  .ratings-table th, .ratings-table td {
    padding: 0.75rem;
    text-align: left;
    border-bottom: 1px solid var(--border-color);
  }

  .ratings-table th {
    background-color: var(--bg-secondary);
    font-weight: bold;
  }

  .ratings-table tr:hover {
    background-color: var(--bg-secondary);
  }

  .edit-button {
    background-color: var(--accent-color);
    color: white;
    border: none;
    padding: 0.4rem 0.8rem;
    border-radius: 4px;
    font-size: 0.8rem;
    cursor: pointer;
  }

  .edit-button:hover {
    background-color: var(--accent-hover);
  }

  .sector-container {
    margin-top: 1rem;
    border: 1px solid var(--border-color);
    border-radius: 8px;
    padding: 1rem;
    max-height: 400px;
    overflow-y: auto;
    background-color: var(--bg-primary);
  }
  
  .sector-header {
    margin-bottom: 1rem;
  }
  
  .sector-header h4 {
    margin: 0 0 0.5rem 0;
  }
  
  .highlighted-sector {
    background-color: var(--highlight-bg);
    padding: 0.5rem;
    border-radius: 4px;
    border-left: 3px solid var(--accent-color);
  }
</style> 