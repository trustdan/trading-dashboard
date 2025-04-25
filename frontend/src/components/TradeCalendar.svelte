<script>
  import { onMount } from 'svelte';
  import { 
    SaveTrade, 
    GetAllTrades, 
    GetTradesByDateRange,
    DeleteTrade
  } from '../../wailsjs/go/main/App.js';
  import { models } from '../../wailsjs/go/models'; // Remove .js extension to allow TypeScript resolution

  // Trade form state - plain object for form binding
  let tradeData = {
    id: 0,
    entryDate: new Date().toISOString().split('T')[0],
    ticker: '',
    sector: '',
    entryPrice: '',
    notes: '',
    // Options-specific fields
    strategyType: '',
    spreadType: '',
    expirationDate: '',
    direction: '' // 'bullish', 'bearish', 'neutral'
  };

  // UI state
  let saving = false;
  let message = '';
  let messageType = 'info';
  let allTrades = [];
  let filteredTrades = [];
  let searchTerm = '';
  let searchField = 'ticker'; // 'ticker', 'sector', 'date', 'strategyType'
  let showStrategyInfo = false;
  let selectedStrategy = '';

  // Calendar state
  let startDate = new Date();
  let weeksToShow = 8;
  let calendarWeeks = [];
  let weeklyTrades = [];
  let selectedWeek = null;
  let showTradeForm = false;
  let editingTrade = null;
  let errorMessage = '';
  let sectors = [];
  let loading = true;
  
  // Sectors for display
  const sectorsForDisplay = [
    'Technology',
    'Healthcare',
    'Financial',
    'Consumer Defensive',
    'Consumer Cyclical',
    'Industrial',
    'Energy',
    'Materials',
    'Utilities',
    'Real Estate',
    'Communication',
    'Others'
  ];

  // Strategy types and categories
  const strategyCategories = [
    {
      name: "Basic Spreads",
      color: "#3498db", // Blue
      strategies: [
        {
          type: "Long Call",
          direction: "bullish",
          description: "Buy a call option, profit from upward price movement",
        },
        {
          type: "Long Put",
          direction: "bearish",
          description: "Buy a put option, profit from downward price movement",
        },
        {
          type: "Short Call",
          direction: "bearish",
          description: "Sell a call option, profit from downward or sideways price movement",
        },
        {
          type: "Short Put",
          direction: "bullish",
          description: "Sell a put option, profit from upward or sideways price movement",
        }
      ]
    },
    {
      name: "Vertical Spreads",
      color: "#9b59b6", // Purple
      strategies: [
        {
          type: "Bull Call Spread",
          direction: "bullish",
          description: "Buy lower strike call, sell higher strike call, same expiration. Also known as Call Debit Spread",
        },
        {
          type: "Bear Call Spread",
          direction: "bearish",
          description: "Sell lower strike call, buy higher strike call, same expiration. Also known as Call Credit Spread",
        },
        {
          type: "Bull Put Spread",
          direction: "bullish",
          description: "Sell higher strike put, buy lower strike put, same expiration. Also known as Put Credit Spread",
        },
        {
          type: "Bear Put Spread",
          direction: "bearish",
          description: "Buy higher strike put, sell lower strike put, same expiration. Also known as Put Debit Spread",
        }
      ]
    },
    {
      name: "Calendar/Horizontal Spreads",
      color: "#2ecc71", // Green
      strategies: [
        {
          type: "Long Calendar Call Spread",
          direction: "neutral",
          description: "Sell near-term call, buy longer-term call, same strike",
        },
        {
          type: "Long Calendar Put Spread",
          direction: "neutral",
          description: "Sell near-term put, buy longer-term put, same strike",
        }
      ]
    },
    {
      name: "Diagonal Spreads",
      color: "#f39c12", // Orange
      strategies: [
        {
          type: "Diagonal Call Spread Up",
          direction: "bullish",
          description: "Buy longer-term lower strike call, sell shorter-term higher strike call",
        },
        {
          type: "Diagonal Call Spread Down",
          direction: "bearish",
          description: "Buy longer-term higher strike call, sell shorter-term lower strike call",
        },
        {
          type: "Diagonal Put Spread Up",
          direction: "bearish",
          description: "Buy longer-term lower strike put, sell shorter-term higher strike put",
        },
        {
          type: "Diagonal Put Spread Down",
          direction: "bullish",
          description: "Buy longer-term higher strike put, sell shorter-term lower strike put",
        }
      ]
    },
    {
      name: "Butterfly Spreads",
      color: "#e74c3c", // Red
      strategies: [
        {
          type: "Long Call Butterfly",
          direction: "neutral",
          description: "Buy 1 lower strike call, sell 2 middle strike calls, buy 1 higher strike call. Strikes evenly spaced, same expiration",
        },
        {
          type: "Long Put Butterfly",
          direction: "neutral",
          description: "Buy 1 lower strike put, sell 2 middle strike puts, buy 1 higher strike put. Strikes evenly spaced, same expiration",
        },
        {
          type: "Broken Wing Butterfly Up",
          direction: "bullish",
          description: "Like standard butterfly but with wider spread between middle and upper strikes",
        },
        {
          type: "Broken Wing Butterfly Down",
          direction: "bearish",
          description: "Like standard butterfly but with wider spread between lower and middle strikes",
        }
      ]
    },
    {
      name: "Iron Condors/Butterflies",
      color: "#1abc9c", // Teal
      strategies: [
        {
          type: "Iron Condor",
          direction: "neutral",
          description: "Sell OTM put, buy further OTM put, sell OTM call, buy further OTM call. All options same expiration",
        },
        {
          type: "Iron Butterfly",
          direction: "neutral",
          description: "Buy OTM put, sell ATM put, sell ATM call, buy OTM call. All options same expiration",
        }
      ]
    },
    {
      name: "Ratio Spreads",
      color: "#d35400", // Brown
      strategies: [
        {
          type: "Call Ratio Backspread",
          direction: "bullish",
          description: "Buy more calls at higher strike than selling at lower strike (e.g., 2:1 ratio)",
        },
        {
          type: "Put Ratio Backspread",
          direction: "bearish",
          description: "Buy more puts at lower strike than selling at higher strike (e.g., 2:1 ratio)",
        },
        {
          type: "Call Ratio Spread",
          direction: "bearish",
          description: "Sell more calls at higher strike than buying at lower strike (e.g., 1:2 ratio)",
        },
        {
          type: "Put Ratio Spread",
          direction: "bullish",
          description: "Sell more puts at lower strike than buying at higher strike (e.g., 1:2 ratio)",
        }
      ]
    }
  ];

  // Flatten strategy types for easy lookup
  let allStrategies = [];
  strategyCategories.forEach(category => {
    category.strategies.forEach(strategy => {
      allStrategies.push({
        ...strategy,
        category: category.name,
        categoryColor: category.color
      });
    });
  });

  // Get color for strategy type
  function getStrategyColor(strategyType) {
    const strategy = allStrategies.find(s => s.type === strategyType);
    return strategy ? strategy.categoryColor : '#95a5a6';
  }

  // Get sector color by name
  const sectorColors = {
    'Technology': '#3498db',
    'Healthcare': '#2ecc71',
    'Financial': '#9b59b6',
    'Consumer Defensive': '#e67e22',
    'Consumer Cyclical': '#e74c3c',
    'Industrial': '#f39c12',
    'Energy': '#d35400',
    'Materials': '#1abc9c',
    'Utilities': '#34495e',
    'Real Estate': '#8e44ad',
    'Communication': '#9b59b6',
    'Others': '#7f8c8d'
  };

  // Direction colors
  const directionColors = {
    'bullish': '#2ecc71',
    'bearish': '#e74c3c',
    'neutral': '#3498db'
  };

  onMount(async () => {
    // Initialize sectors array from sectorsForDisplay
    sectors = [...sectorsForDisplay];
    
    // Generate calendar weeks
    generateCalendarWeeks();
    
    // Load all trades (or use mock data if backend unavailable)
    try {
      await loadAllTrades();
    } catch (err) {
      console.error('Backend unavailable, using mock data:', err);
      useMockData();
    }
  });

  // Generate calendar weeks for the next N weeks
  function generateCalendarWeeks() {
    calendarWeeks = [];
    // Ensure startDate is set to the beginning of the current week (Sunday)
    const today = new Date();
    const dayOfWeek = today.getDay(); // 0 = Sunday, 1 = Monday, etc.
    
    // Adjust to the previous Sunday to start the week
    const startOfWeek = new Date(today);
    startOfWeek.setDate(today.getDate() - dayOfWeek);
    
    // Update startDate to the beginning of the current week
    startDate = new Date(startOfWeek);
    
    console.log("Calendar start date:", startDate.toISOString());
    
    const currentDate = new Date(startDate);
    
    for (let i = 0; i < weeksToShow; i++) {
      // Create a week object
      const weekStart = new Date(currentDate);
      const weekEnd = new Date(currentDate);
      weekEnd.setDate(weekEnd.getDate() + 6); // Include the full week
      
      const week = {
        startDate: new Date(weekStart),
        endDate: new Date(weekEnd),
        weekNumber: i + 1,
        // Get Friday of the week (for options expiration)
        expirationDate: getNextFriday(new Date(currentDate)),
        sectors: {}
      };
      
      // Log each week's start and end dates for debugging
      console.log(`Week ${i+1}: ${week.startDate.toISOString()} to ${week.endDate.toISOString()}, expiry: ${week.expirationDate.toISOString()}`);
      
      // Initialize trades for each sector
      sectorsForDisplay.forEach(sector => {
        week.sectors[sector] = [];
      });
      
      // Move to next week
      currentDate.setDate(currentDate.getDate() + 7);
      calendarWeeks.push(week);
    }
    
    console.log("Generated calendar weeks:", calendarWeeks);
    
    // Map existing trades to weeks
    if (allTrades.length > 0) {
      mapTradesToWeeks();
    }
  }

  // Helper function to get next Friday from a date
  function getNextFriday(date) {
    const dayOfWeek = date.getDay();
    const daysUntilFriday = (5 - dayOfWeek + 7) % 7;
    const friday = new Date(date);
    friday.setDate(date.getDate() + daysUntilFriday);
    return friday;
  }

  // Map trades to their respective weeks and sectors
  function mapTradesToWeeks() {
    weeklyTrades = [];
    
    // Initialize the weeklyTrades array with empty arrays for each week and sector
    for (let i = 0; i < calendarWeeks.length; i++) {
      weeklyTrades[i] = {};
      for (const sector of sectorsForDisplay) {
        weeklyTrades[i][sector] = [];
      }
    }
    
    console.log("Mapping trades to weeks. Total trades:", allTrades.length);
    
    // Map trades to their respective weeks and sectors
    for (const trade of allTrades) {
      const entryDate = new Date(trade.entryDate);
      const expiryDate = new Date(trade.expirationDate);
      
      // Find the entry week and expiry week indices
      const entryWeekIndex = calendarWeeks.findIndex(week => {
        const weekStart = new Date(week.startDate);
        const weekEnd = new Date(week.endDate);
        return entryDate >= weekStart && entryDate <= weekEnd;
      });
      
      const expiryWeekIndex = calendarWeeks.findIndex(week => {
        const weekStart = new Date(week.startDate);
        const weekEnd = new Date(week.endDate);
        return expiryDate >= weekStart && expiryDate <= weekEnd;
      });
      
      console.log(`Trade ${trade.ticker}: Entry week ${entryWeekIndex}, Expiry week ${expiryWeekIndex}`);
      
      // If either week is not found, skip this trade
      if (entryWeekIndex === -1 || expiryWeekIndex === -1) {
        console.log(`Skipping trade ${trade.ticker} - outside calendar range`);
        continue;
      }
      
      // Add the trade to each week between entry and expiry (inclusive)
      for (let i = entryWeekIndex; i <= expiryWeekIndex; i++) {
        if (weeklyTrades[i] && weeklyTrades[i][trade.sector]) {
          // Create a clone of the trade object with additional display properties
          const tradeClone = { 
            ...trade,
            isEntryWeek: i === entryWeekIndex,
            isExpiryWeek: i === expiryWeekIndex,
            isPartialWeek: i > entryWeekIndex && i < expiryWeekIndex,
            weekSpan: expiryWeekIndex - entryWeekIndex + 1,
            weekPosition: i - entryWeekIndex + 1
            // Note: We're no longer copying convertValues to avoid issues with plain objects
          };
          
          weeklyTrades[i][trade.sector].push(tradeClone);
        }
      }
    }
    
    console.log("Weekly trades mapping complete", weeklyTrades);
  }

  // Use mock data when backend is unavailable
  function useMockData() {
    // Sample mock data - ensure dates are properly formatted
    const today = new Date();
    
    // First week options - ensure these are YYYY-MM-DD strings
    const friday1 = getNextFriday(today).toISOString().split('T')[0];
    // Second week options  
    const nextWeek = new Date(today);
    nextWeek.setDate(today.getDate() + 7);
    const friday2 = getNextFriday(nextWeek).toISOString().split('T')[0];
    // Third week options
    const twoWeeksLater = new Date(today);
    twoWeeksLater.setDate(today.getDate() + 14);
    const friday3 = getNextFriday(twoWeeksLater).toISOString().split('T')[0];
    
    console.log("Mock data expiration dates:", {
      friday1,
      friday2,
      friday3
    });
    
    // Create mock trades as plain objects first
    const mockTradesData = [
      {
        id: 1,
        entryDate: today.toISOString().split('T')[0],
        expirationDate: friday1,
        ticker: 'AAPL',
        sector: 'Technology',
        entryPrice: '15.50',
        strategyType: 'Long Call',
        direction: 'bullish',
        notes: 'DEMO - Bullish on earnings'
      },
      {
        id: 2,
        entryDate: today.toISOString().split('T')[0],
        expirationDate: friday1,
        ticker: 'MSFT',
        sector: 'Technology',
        entryPrice: '8.25',
        strategyType: 'Bull Call Spread',
        direction: 'bullish',
        notes: 'DEMO - Moderate upside expected'
      },
      {
        id: 3,
        entryDate: today.toISOString().split('T')[0],
        expirationDate: friday2,
        ticker: 'AMZN',
        sector: 'Consumer Cyclical',
        entryPrice: '12.75',
        strategyType: 'Iron Condor',
        direction: 'neutral',
        notes: 'DEMO - Expecting sideways movement'
      },
      {
        id: 4,
        entryDate: today.toISOString().split('T')[0],
        expirationDate: friday3,
        ticker: 'JPM',
        sector: 'Financial',
        entryPrice: '5.50',
        strategyType: 'Long Put',
        direction: 'bearish',
        notes: 'DEMO - Bearish on banking sector'
      },
      {
        id: 5,
        entryDate: today.toISOString().split('T')[0],
        expirationDate: friday2,
        ticker: 'KO',
        sector: 'Consumer Defensive',
        entryPrice: '3.25',
        strategyType: 'Long Calendar Put Spread',
        direction: 'neutral',
        notes: 'DEMO - Stable stock during market volatility'
      }
    ];
    
    // Convert each mock trade to a proper Trade instance
    allTrades = mockTradesData.map(tradeData => {
      // Keep entryPrice as string for the form binding
      return tradeData;
    });
    
    console.log("Using mock data:", JSON.stringify(allTrades));
    
    // Force reactivity
    allTrades = [...allTrades];
    filterTrades();
    mapTradesToWeeks();
    
    message = 'Connected to demo mode. Backend unavailable.';
    messageType = 'info';
  }

  // Load all trades
  async function loadAllTrades() {
    try {
      allTrades = await GetAllTrades();
      filterTrades();
      mapTradesToWeeks();
    } catch (err) {
      console.error('Error fetching trades:', err);
      message = 'Failed to load trades from backend';
      messageType = 'error';
      throw err;
    }
  }

  // Filter trades based on search criteria
  function filterTrades() {
    if (!searchTerm) {
      filteredTrades = [...allTrades];
      return;
    }
    
    const term = searchTerm.toLowerCase();
    
    // Filter based on selected field
    switch (searchField) {
      case 'ticker':
        filteredTrades = allTrades.filter(trade => 
          trade.ticker.toLowerCase().includes(term)
        );
        break;
      case 'sector':
        filteredTrades = allTrades.filter(trade => 
          trade.sector.toLowerCase().includes(term)
        );
        break;
      case 'strategyType':
        filteredTrades = allTrades.filter(trade => 
          (trade.strategyType || '').toLowerCase().includes(term)
        );
        break;
      case 'date':
        filteredTrades = allTrades.filter(trade => 
          new Date(trade.entryDate).toLocaleDateString().includes(term) ||
          (trade.expirationDate && new Date(trade.expirationDate).toLocaleDateString().includes(term))
        );
        break;
      default:
        filteredTrades = [...allTrades];
    }
  }

  // Handle search input changes
  function handleSearch() {
    filterTrades();
  }

  // Show strategy information
  function showStrategyDetails(strategy) {
    selectedStrategy = strategy;
    showStrategyInfo = true;
  }

  // Hide strategy information
  function hideStrategyInfo() {
    showStrategyInfo = false;
  }

  // Update direction when strategy is selected
  function updateDirection() {
    if (tradeData.strategyType) {
      const strategy = allStrategies.find(s => s.type === tradeData.strategyType);
      if (strategy) {
        tradeData.direction = strategy.direction;
      }
    }
  }

  // Save the current trade
  async function saveTrade() {
    // Basic validation
    if (!tradeData.ticker) {
      message = 'Please enter a ticker symbol';
      messageType = 'error';
      setTimeout(() => { message = ''; }, 3000);
      return;
    }
    if (!tradeData.sector) {
      message = 'Please select a sector';
      messageType = 'error';
      setTimeout(() => { message = ''; }, 3000);
      return;
    }
    if (!tradeData.entryPrice) {
      message = 'Please enter an entry price';
      messageType = 'error';
      setTimeout(() => { message = ''; }, 3000);
      return;
    }
    if (!tradeData.strategyType) {
      message = 'Please select an options strategy';
      messageType = 'error';
      setTimeout(() => { message = ''; }, 3000);
      return;
    }
    if (!tradeData.expirationDate) {
      message = 'Please select an expiration date';
      messageType = 'error';
      setTimeout(() => { message = ''; }, 3000);
      return;
    }

    saving = true;
    message = '';
    let savedSuccessfully = false;
    
    // Create a models.Trade instance from the form data
    const tradeToSave = models.Trade.createFrom({
      ...tradeData,
      entryPrice: parseFloat(tradeData.entryPrice) // Ensure entryPrice is a number
    });

    try {
      let result;
      try {
        result = await SaveTrade(tradeToSave); // Pass the models.Trade instance
        // Backend save successful
        if (tradeToSave.id === 0) {
          // Assuming backend returns the saved object with ID
          allTrades.push(models.Trade.createFrom(result)); // Store as Trade instance
        } else {
          const index = allTrades.findIndex(t => t.id === tradeToSave.id);
          if (index !== -1) {
            allTrades[index] = models.Trade.createFrom(result); // Update with Trade instance
          }
        }
        savedSuccessfully = true;
      } catch (err) {
        console.warn('Backend unavailable during save, using demo mode:', err);
        // Mock save in demo mode
        if (tradeToSave.id === 0) {
          tradeToSave.id = allTrades.length > 0 ? Math.max(...allTrades.map(t => t.id)) + 1 : 1;
          // Ensure dates are correct, createFrom handles this
          const savedTrade = models.Trade.createFrom(tradeToSave);
          allTrades.push(savedTrade);
          console.log("Added new trade to allTrades (demo):", JSON.stringify(savedTrade));
        } else {
          const index = allTrades.findIndex(t => t.id === tradeToSave.id);
          if (index !== -1) {
            allTrades[index] = models.Trade.createFrom(tradeToSave);
            console.log("Updated existing trade in allTrades (demo) at index", index);
          }
        }
        result = models.Trade.createFrom(tradeToSave); // Use the local trade object as the result
        savedSuccessfully = true;
      }
      
      if (savedSuccessfully) {
        // Force reactivity
        allTrades = [...allTrades]; 
        console.log("After save, allTrades length:", allTrades.length);
        
        message = 'Trade saved successfully';
        messageType = 'success';
        
        // Reset form for new entry
        resetTradeForm(); // Use the reset function
        
        // Update calendar and table
        filterTrades(); // Ensure filteredTrades is updated
        mapTradesToWeeks(); // Remap trades including the new/updated one
      } else {
        message = 'Failed to save trade';
        messageType = 'error';
      }
    } catch (err) {
      console.error('Error during post-save update:', err);
      message = 'Failed to update view after saving trade';
      messageType = 'error';
    } finally {
      saving = false;
      setTimeout(() => { message = ''; }, 3000);
    }
  }

  // Delete a trade
  async function deleteTrade(tradeId) {
    if (!confirm('Are you sure you want to delete this trade?')) {
      return;
    }

    message = '';
    let deletedSuccessfully = false;
    
    try {
      // Try backend delete
      try {
        await DeleteTrade(tradeId);
        deletedSuccessfully = true;
      } catch(err) {
        console.warn('Backend unavailable during delete, using demo mode:', err);
        // If backend fails, assume demo mode allows local delete
        deletedSuccessfully = true;
      }
      
      if (deletedSuccessfully) {
        allTrades = allTrades.filter(t => t.id !== tradeId);
        allTrades = [...allTrades]; // Trigger reactivity
        
        filterTrades();
        mapTradesToWeeks();
        
        message = 'Trade deleted successfully';
        messageType = 'success';
      } else {
        message = 'Failed to delete trade';
        messageType = 'error';
      }
    } catch (err) {
      console.error('Error deleting trade:', err);
      message = 'Failed to delete trade';
      messageType = 'error';
    } finally {
      setTimeout(() => {
        message = '';
      }, 3000);
    }
  }

  // Edit an existing trade
  function editTrade(existingTrade) {
    // Use tradeData for form binding
    tradeData = { 
      ...existingTrade,
      // Ensure dates are in YYYY-MM-DD format for input type="date"
      entryDate: formatDate(existingTrade.entryDate),
      expirationDate: existingTrade.expirationDate ? formatDate(existingTrade.expirationDate) : '',
      entryPrice: existingTrade.entryPrice.toString() // Convert entry price back to string for input
    };
    
    // Scroll to trade form
    document.querySelector('.trade-form').scrollIntoView({ 
      behavior: 'smooth' 
    });
  }

  // Reset the trade form
  function resetTradeForm() {
    tradeData = {
      id: 0,
      entryDate: new Date().toISOString().split('T')[0],
      ticker: '',
      sector: '',
      entryPrice: '',
      notes: '',
      strategyType: '',
      spreadType: '',
      expirationDate: '',
      direction: ''
    };
  }

  // Get a sector color from the map or fallback to a default
  function getSectorColor(sector) {
    return sectorColors[sector] || '#95a5a6';
  }

  // Get a direction color
  function getDirectionColor(direction) {
    return directionColors[direction] || '#95a5a6';
  }

  // Format a date as YYYY-MM-DD
  function formatDate(date) {
    const d = new Date(date);
    return d.toISOString().split('T')[0];
  }

  // Format a date for display
  function formatDisplayDate(date) {
    return new Date(date).toLocaleDateString(undefined, { 
      month: 'short', 
      day: 'numeric',
      year: '2-digit'
    });
  }

  function renderTradeItem(trade, weekIndex) {
    const classes = [
      'trade-pill',
      trade.isEntryWeek ? 'entry-week' : '',
      trade.isExpiryWeek ? 'expiry-week' : '',
      trade.isPartialWeek ? 'isPartialWeek' : '',
      trade.isEntryWeek && trade.isExpiryWeek ? '' : 'multi-week'
    ].filter(Boolean).join(' ');
    
    return `
      <div class="${classes}" 
           data-trade-id="${trade.id}" 
           style="--accent-color: ${getDirectionColor(trade.direction)}">
        <div class="trade-ticker">${trade.ticker}</div>
        <div class="trade-strategy">${trade.strategyType} ${trade.spreadType}</div>
      </div>
    `;
  }
</script>

<div class="trade-calendar">
  <h2>Options Trading Calendar</h2>
  <p class="dashboard-description">
    Track and visualize your options trades across sectors and expiration weeks.
  </p>

  {#if message}
    <div class="message {messageType}">
      {message}
    </div>
  {/if}
  
  <div class="strategy-guide">
    <h3>Strategy Guide</h3>
    <div class="strategy-colors">
      {#each strategyCategories as category}
        <div class="strategy-color-item">
          <div class="color-box" style="background-color: {category.color}"></div>
          <div class="color-label">{category.name}</div>
        </div>
      {/each}
    </div>
  </div>

  <div class="calendar-container">
    <div class="calendar">
      <!-- Header with week numbers and expiration dates -->
      <div class="calendar-header">
        <div class="sector-header">Sector</div>
        {#each calendarWeeks as week}
          <div class="week-header">
            <div>Week {week.weekNumber}</div>
            <div class="expiration-date">Exp: {formatDisplayDate(week.expirationDate)}</div>
          </div>
        {/each}
      </div>

      <!-- Grid with sectors and trades -->
      {#each sectorsForDisplay as sector}
        <div class="calendar-row">
          <div class="sector-cell" style="background-color: {getSectorColor(sector)}30;">
            {sector}
          </div>
          
          {#each calendarWeeks as week, weekIndex}
            <div class="week-cell">
              {#if weeklyTrades[weekIndex] && weeklyTrades[weekIndex][sector]}
                {#each weeklyTrades[weekIndex][sector] as tradeItem}
                  <div 
                    class="trade-pill {tradeItem.isEntryWeek ? 'entry-week' : ''} {tradeItem.isExpiryWeek ? 'expiry-week' : ''} {tradeItem.isPartialWeek ? 'partial-week' : ''}"
                    style="background-color: {getStrategyColor(tradeItem.strategyType)}"
                    on:click={() => editTrade(tradeItem)} 
                  >
                    <div class="trade-pill-ticker">{tradeItem.ticker}</div>
                    <div class="trade-pill-strategy">{tradeItem.strategyType} {tradeItem.spreadType || ''}</div>
                    <div class="trade-direction-indicator" style="background-color: {getDirectionColor(tradeItem.direction)}"></div>
                  </div>
                {/each}
              {/if}
            </div>
          {/each}
        </div>
      {/each}
    </div>
  </div>

  <div class="trade-form">
    <h3>{tradeData.id ? 'Edit Options Trade' : 'Add New Options Trade'}</h3>
    
    <div class="form-container">
      <div class="form-row">
        <div class="form-group">
          <label for="entry-date">Entry Date:</label>
          <input 
            type="date" 
            id="entry-date" 
            bind:value={tradeData.entryDate}
          />
        </div>

        <div class="form-group">
          <label for="expiration-date">Expiration Date:</label>
          <input 
            type="date" 
            id="expiration-date" 
            bind:value={tradeData.expirationDate}
          />
        </div>

        <div class="form-group">
          <label for="ticker">Ticker:</label>
          <input 
            type="text" 
            id="ticker" 
            bind:value={tradeData.ticker}
            placeholder="e.g., AAPL"
            maxlength="10"
          />
        </div>

        <div class="form-group">
          <label for="sector">Sector:</label>
          <select id="sector" bind:value={tradeData.sector}>
            <option value="">Select a sector...</option>
            {#each sectorsForDisplay as sector}
              <option value={sector}>{sector}</option>
            {/each}
          </select>
        </div>

        <div class="form-group">
          <label for="entry-price">Entry Price:</label>
          <input 
            type="number" 
            id="entry-price" 
            bind:value={tradeData.entryPrice}
            step="0.01"
            min="0"
            placeholder="0.00"
          />
        </div>
      </div>

      <div class="form-row">
        <div class="form-group strategy-group">
          <label for="strategy-type">Options Strategy:</label>
          <select id="strategy-type" bind:value={tradeData.strategyType} on:change={updateDirection}>
            <option value="">Select a strategy...</option>
            {#each strategyCategories as category}
              <optgroup label={category.name}>
                {#each category.strategies as strategy}
                  <option value={strategy.type}>{strategy.type}</option>
                {/each}
              </optgroup>
            {/each}
          </select>
          
          {#if tradeData.strategyType}
            <div class="strategy-info">
              {#if tradeData.direction}
                <span class="direction-badge" style="background-color: {getDirectionColor(tradeData.direction)}">
                  {tradeData.direction.charAt(0).toUpperCase() + tradeData.direction.slice(1)}
                </span>
              {/if}
              <button type="button" class="info-button" on:click={() => showStrategyDetails(tradeData.strategyType)}>
                Strategy Info
              </button>
            </div>
          {/if}
        </div>
      </div>

      <div class="form-row">
        <div class="form-group notes-group">
          <label for="notes">Trade Notes:</label>
          <textarea 
            id="notes" 
            bind:value={tradeData.notes}
            placeholder="Add any relevant trade notes, strategy details, or observations..."
            rows="3"
          ></textarea>
        </div>
      </div>

      <div class="form-actions">
        <button class="reset-button" on:click={resetTradeForm}>Reset</button>
        <button class="save-button" on:click={saveTrade} disabled={saving}>
          {saving ? 'Saving...' : (tradeData.id ? 'Update Trade' : 'Save Trade')}
        </button>
      </div>
    </div>
  </div>

  <div class="trade-history">
    <h3>Options Trade History</h3>
    
    <div class="search-container">
      <div class="search-field">
        <input 
          type="text" 
          placeholder="Search trades..." 
          bind:value={searchTerm}
          on:input={handleSearch}
        />
      </div>
      <div class="search-options">
        <label>
          <input type="radio" bind:group={searchField} value="ticker" on:change={handleSearch}>
          Ticker
        </label>
        <label>
          <input type="radio" bind:group={searchField} value="sector" on:change={handleSearch}>
          Sector
        </label>
        <label>
          <input type="radio" bind:group={searchField} value="strategyType" on:change={handleSearch}>
          Strategy
        </label>
        <label>
          <input type="radio" bind:group={searchField} value="date" on:change={handleSearch}>
          Date
        </label>
      </div>
    </div>
    
    {#if filteredTrades.length === 0}
      <p class="no-trades">No trades found matching your criteria.</p>
    {:else}
      <div class="trades-table-container">
        <table class="trades-table">
          <thead>
            <tr>
              <th>Entry Date</th>
              <th>Expiration</th>
              <th>Ticker</th>
              <th>Sector</th>
              <th>Strategy</th>
              <th>Direction</th>
              <th>Entry Price</th>
              <th>Notes</th>
              <th>Action</th>
            </tr>
          </thead>
          <tbody>
            {#each filteredTrades as item}
              <tr>
                <td>{new Date(item.entryDate).toLocaleDateString()}</td>
                <td>{item.expirationDate ? new Date(item.expirationDate).toLocaleDateString() : '-'}</td>
                <td>{item.ticker}</td>
                <td>
                  <span class="sector-badge" style="background-color: {getSectorColor(item.sector)}">
                    {item.sector}
                  </span>
                </td>
                <td>{item.strategyType || '-'}</td>
                <td>
                  {#if item.direction}
                    <span class="direction-badge" style="background-color: {getDirectionColor(item.direction)}">
                      {item.direction.charAt(0).toUpperCase() + item.direction.slice(1)}
                    </span>
                  {:else}
                    -
                  {/if}
                </td>
                <td>${item.entryPrice ? parseFloat(item.entryPrice).toFixed(2) : '0.00'}</td>
                <td class="notes-cell">{item.notes}</td>
                <td class="action-cell">
                  <button class="edit-button" on:click={() => editTrade(item)}>Edit</button>
                  <button class="delete-button" on:click={() => deleteTrade(item.id)}>Delete</button>
                </td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>
    {/if}
  </div>

  <!-- Strategy Info Modal -->
  {#if showStrategyInfo && selectedStrategy}
    <div class="modal-overlay" on:click={hideStrategyInfo}>
      <div class="modal-content" on:click|stopPropagation>
        <button class="close-button" on:click={hideStrategyInfo}>×</button>
        
        {#if selectedStrategy}
          {@const strategy = allStrategies.find(s => s.type === selectedStrategy)}
          {#if strategy}
            <h3>{strategy.type}</h3>
            <div class="strategy-category">{strategy.category}</div>
            
            <div class="strategy-direction">
              <span class="direction-badge" style="background-color: {getDirectionColor(strategy.direction)}">
                {strategy.direction.charAt(0).toUpperCase() + strategy.direction.slice(1)}
              </span>
            </div>
            
            <div class="strategy-description">
              {strategy.description}
            </div>
          {/if}
        {/if}
      </div>
    </div>
  {/if}
</div>

<style>
  .trade-calendar {
    max-width: 1200px;
    margin: 0 auto;
  }

  h2 {
    color: var(--text-primary);
    border-bottom: 2px solid var(--accent-color, #3498db);
    padding-bottom: 0.5rem;
  }

  h3 {
    color: var(--text-primary);
    margin: 0;
  }

  .dashboard-description {
    color: var(--text-secondary);
    margin-bottom: 2rem;
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

  /* Calendar Styles */
  .calendar-container {
    background-color: var(--bg-primary);
    border-radius: 8px;
    box-shadow: 0 2px 10px var(--shadow-color);
    padding: 1.5rem;
    margin-bottom: 2rem;
    overflow: auto;
  }

  .calendar {
    border: 1px solid var(--border-color);
    border-radius: 4px;
    min-width: 900px; /* To handle overflow on smaller screens */
  }

  .calendar-header {
    display: grid;
    grid-template-columns: 180px repeat(8, 1fr);
    background-color: var(--bg-secondary);
    border-bottom: 1px solid var(--border-color);
    position: sticky;
    top: 0;
    z-index: 1;
  }

  .sector-header, .week-header {
    padding: 0.75rem;
    text-align: center;
    font-weight: bold;
    color: var(--text-primary);
  }

  .sector-header {
    text-align: left;
    border-right: 1px solid var(--border-color);
  }

  .expiration-date {
    font-size: 0.8rem;
    color: var(--text-secondary);
    margin-top: 0.25rem;
  }

  .calendar-row {
    display: grid;
    grid-template-columns: 180px repeat(8, 1fr);
    border-bottom: 1px solid var(--border-color);
  }

  .calendar-row:last-child {
    border-bottom: none;
  }

  .sector-cell {
    padding: 0.75rem;
    font-weight: bold;
    border-right: 1px solid var(--border-color);
    display: flex;
    align-items: center;
    color: var(--text-primary);
  }

  .week-cell {
    min-height: 60px;
    padding: 0.5rem;
    border-right: 1px solid var(--border-color);
    background-color: var(--bg-primary);
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
  }

  .week-cell:last-child {
    border-right: none;
  }

  .trade-pill {
    position: relative;
    background-color: var(--accent-color, #3498db);
    color: white;
    font-size: 0.8rem;
    padding: 0.4rem 0.5rem;
    padding-left: 0.7rem;
    border-radius: 4px;
    cursor: pointer;
    transition: opacity 0.3s;
    overflow: hidden;
    display: flex;
    flex-direction: column;
  }

  .trade-pill.multi-week {
    position: relative;
  }

  /* Add visual markers for entry and expiry weeks */
  .trade-pill.entry-week:before {
    content: '';
    position: absolute;
    left: 0;
    top: 0;
    bottom: 0;
    width: 4px;
    background-color: #2ecc71; /* Green for entry */
  }
  
  .trade-pill.expiry-week:after {
    content: '';
    position: absolute;
    right: 0;
    top: 0;
    bottom: 0;
    width: 4px;
    background-color: #e74c3c; /* Red for expiry */
  }

  .trade-pill.entry-week.expiry-week:before {
    left: 0;
    width: 2px;
  }
  
  .trade-pill.entry-week.expiry-week:after {
    right: 0;
    width: 2px;
  }

  /* Add styles for trades that span multiple weeks */
  .trade-pill.partial-week {
    border-radius: 0;
    margin-left: -1px;
    margin-right: -1px;
  }

  .trade-pill.entry-week:not(.expiry-week) {
    border-top-right-radius: 0;
    border-bottom-right-radius: 0;
    margin-right: -1px;
  }

  .trade-pill.expiry-week:not(.entry-week) {
    border-top-left-radius: 0;
    border-bottom-left-radius: 0;
    margin-left: -1px;
  }

  .trade-pill:hover {
    opacity: 0.9;
    transform: scale(1.02);
    z-index: 10;
  }

  .trade-pill-ticker {
    font-weight: bold;
  }

  .trade-pill-strategy {
    font-size: 0.7rem;
    opacity: 0.9;
  }

  /* Trade Form Styles */
  .trade-form {
    background-color: var(--bg-primary);
    border-radius: 8px;
    box-shadow: 0 2px 10px var(--shadow-color);
    padding: 1.5rem;
    margin-bottom: 2rem;
  }

  .form-container {
    margin-top: 1rem;
  }

  .form-row {
    display: flex;
    gap: 1rem;
    margin-bottom: 1rem;
  }

  .form-group {
    flex: 1;
  }

  .form-group label {
    display: block;
    margin-bottom: 0.5rem;
    font-weight: bold;
    color: var(--text-primary);
  }

  .form-group input,
  .form-group select,
  .form-group textarea {
    width: 100%;
    padding: 0.5rem;
    border: 1px solid var(--border-color);
    border-radius: 4px;
    background-color: var(--bg-primary);
    color: var(--text-primary);
  }

  .strategy-group {
    flex: 2;
  }

  .strategy-info {
    display: flex;
    align-items: center;
    margin-top: 0.5rem;
    gap: 0.5rem;
  }

  .direction-badge {
    display: inline-block;
    color: white;
    padding: 0.25rem 0.5rem;
    border-radius: 4px;
    font-size: 0.8rem;
  }

  .info-button {
    background-color: var(--neutral-color, #95a5a6);
    color: white;
    border: none;
    padding: 0.25rem 0.5rem;
    border-radius: 4px;
    font-size: 0.8rem;
    cursor: pointer;
  }

  .info-button:hover {
    background-color: var(--neutral-hover, #7f8c8d);
  }

  .notes-group {
    flex: 1;
  }

  .form-actions {
    display: flex;
    justify-content: flex-end;
    gap: 1rem;
    margin-top: 1rem;
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

  .reset-button {
    background-color: var(--neutral-color, #95a5a6);
    color: white;
    border: none;
    padding: 0.8rem 2rem;
    border-radius: 4px;
    font-size: 1rem;
    cursor: pointer;
    transition: background-color 0.3s;
  }

  .reset-button:hover {
    background-color: var(--neutral-hover, #7f8c8d);
  }

  /* Trade History Styles */
  .trade-history {
    background-color: var(--bg-primary);
    border-radius: 8px;
    box-shadow: 0 2px 10px var(--shadow-color);
    padding: 1.5rem;
  }

  .search-container {
    display: flex;
    margin: 1rem 0;
    gap: 1rem;
    align-items: center;
  }

  .search-field {
    flex: 1;
  }

  .search-field input {
    width: 100%;
    padding: 0.5rem;
    border: 1px solid var(--border-color);
    border-radius: 4px;
    background-color: var(--bg-primary);
    color: var(--text-primary);
  }

  .search-options {
    display: flex;
    gap: 1rem;
    color: var(--text-primary);
  }

  .search-options label {
    display: flex;
    align-items: center;
    gap: 0.25rem;
    cursor: pointer;
  }

  .no-trades {
    text-align: center;
    color: var(--text-secondary);
    font-style: italic;
    margin: 2rem 0;
  }

  .trades-table-container {
    overflow-x: auto;
  }

  .trades-table {
    width: 100%;
    border-collapse: collapse;
  }

  .trades-table th, .trades-table td {
    padding: 0.75rem;
    text-align: left;
    border-bottom: 1px solid var(--border-color);
    color: var(--text-primary);
  }

  .trades-table th {
    background-color: var(--bg-secondary);
    font-weight: bold;
  }

  .trades-table tr:hover {
    background-color: var(--bg-secondary);
  }

  .sector-badge {
    display: inline-block;
    color: white;
    padding: 0.25rem 0.5rem;
    border-radius: 4px;
    font-size: 0.8rem;
  }

  .notes-cell {
    max-width: 300px;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .edit-button {
    background-color: var(--accent-color, #3498db);
    color: white;
    border: none;
    padding: 0.4rem 0.8rem;
    border-radius: 4px;
    font-size: 0.8rem;
    cursor: pointer;
  }

  .edit-button:hover {
    background-color: var(--accent-hover, #2980b9);
  }

  /* Updated Action Cell Styles */
  .action-cell {
    display: flex;
    gap: 0.5rem;
  }

  .delete-button {
    background-color: var(--error-color, #e74c3c);
    color: white;
    border: none;
    padding: 0.4rem 0.8rem;
    border-radius: 4px;
    font-size: 0.8rem;
    cursor: pointer;
  }

  .delete-button:hover {
    background-color: var(--error-hover, #c0392b); /* Need to define --error-hover */
  }

  /* Modal styles */
  .modal-overlay {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background-color: rgba(0, 0, 0, 0.5);
    display: flex;
    justify-content: center;
    align-items: center;
    z-index: 1000;
  }

  .modal-content {
    background-color: var(--bg-primary);
    border-radius: 8px;
    box-shadow: 0 2px 10px var(--shadow-color);
    padding: 1.5rem;
    max-width: 500px;
    width: 100%;
    position: relative;
  }

  .close-button {
    position: absolute;
    top: 0.5rem;
    right: 0.5rem;
    background: none;
    border: none;
    font-size: 1.5rem;
    cursor: pointer;
    color: var(--text-secondary);
  }

  .strategy-category {
    color: var(--text-secondary);
    margin-bottom: 0.5rem;
  }

  .strategy-direction {
    margin-bottom: 1rem;
  }

  .strategy-description {
    line-height: 1.5;
    color: var(--text-primary);
  }

  /* Strategy guide styles */
  .strategy-guide {
    background-color: var(--bg-primary);
    border-radius: 8px;
    box-shadow: 0 2px 10px var(--shadow-color);
    padding: 1rem;
    margin-bottom: 1rem;
  }
  
  .strategy-guide h3 {
    margin-top: 0;
    margin-bottom: 0.5rem;
    font-size: 1rem;
  }
  
  .strategy-colors {
    display: flex;
    flex-wrap: wrap;
    gap: 1rem;
  }
  
  .strategy-color-item {
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }
  
  .color-box {
    width: 16px;
    height: 16px;
    border-radius: 4px;
  }
  
  .color-label {
    font-size: 0.8rem;
    color: var(--text-primary);
  }

  /* Responsive adjustments */
  @media (max-width: 768px) {
    .form-row {
      flex-direction: column;
    }
  }
</style> 