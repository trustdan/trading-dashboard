<script>
  import { GetVersion } from '../wailsjs/go/main/App.js';
  import { onMount } from 'svelte';

  // Components import (will be created next)
  import RiskDashboard from './components/RiskDashboard.svelte';
  import StockRating from './components/StockRating.svelte';
  import TradeCalendar from './components/TradeCalendar.svelte';

  // Navigation state
  let activeTab = 'risk'; // 'risk', 'stock', 'calendar'
  let appVersion = '';
  
  // Theme state
  let darkMode = localStorage.getItem('darkMode') === 'true';
  
  // Update localStorage and body class when theme changes
  $: {
    localStorage.setItem('darkMode', darkMode ? 'true' : 'false');
    if (typeof document !== 'undefined') {
      document.body.classList.toggle('dark-mode', darkMode);
    }
  }
  
  // Toggle theme function
  function toggleTheme() {
    darkMode = !darkMode;
  }

  // On component mount, get app version
  onMount(async () => {
    try {
      appVersion = await GetVersion();
    } catch (err) {
      console.error('Error fetching app version:', err);
    }
  });

  // Handle tab switch
  function switchTab(tab) {
    activeTab = tab;
  }
</script>

<main class={darkMode ? 'dark-theme' : 'light-theme'}>
  <div class="app-container">
    <header class="app-header">
      <h1>Trading Dashboard</h1>
      <div class="header-right">
        <button class="theme-toggle" on:click={toggleTheme}>
          {darkMode ? '☀️' : '🌙'}
        </button>
        <div class="version">v{appVersion}</div>
      </div>
    </header>

    <nav class="app-nav">
      <ul>
        <li class:active={activeTab === 'risk'}>
          <button on:click={() => switchTab('risk')}>Risk Management</button>
        </li>
        <li class:active={activeTab === 'stock'}>
          <button on:click={() => switchTab('stock')}>Stock Rating</button>
        </li>
        <li class:active={activeTab === 'calendar'}>
          <button on:click={() => switchTab('calendar')}>Trade Calendar</button>
        </li>
      </ul>
    </nav>

    <div class="app-content">
      {#if activeTab === 'risk'}
        <RiskDashboard />
      {:else if activeTab === 'stock'}
        <StockRating />
      {:else if activeTab === 'calendar'}
        <TradeCalendar />
      {/if}
    </div>

    <footer class="app-footer">
      <p>Trading Dashboard • Local-First Trading Journal</p>
    </footer>
  </div>
</main>

<style>
  /* Theme Variables - Light Theme (Default) */
  .light-theme {
    --bg-main: #f5f7fa;
    --bg-primary: #ffffff;
    --bg-secondary: #f8f9fa;
    --text-primary: #333;
    --text-secondary: #666;
    --header-bg: #2c3e50;
    --header-text: #ffffff;
    --nav-bg: #34495e;
    --nav-text: #ffffff;
    --nav-active: #1abc9c;
    --border-color: #ddd;
    --shadow-color: rgba(0, 0, 0, 0.1);
    --accent-color: #3498db;
    --accent-hover: #2980b9;
    --neutral-color: #95a5a6;
    --neutral-hover: #7f8c8d;
    --success-color: #27ae60;
    --success-bg: #d5f5e3;
    --warning-color: #f39c12;
    --error-color: #e74c3c;
    --error-bg: #fadbd8;
    --info-bg: #d6eaf8;
  }
  
  /* Dark Theme */
  .dark-theme {
    --bg-main: #1a1a2e;
    --bg-primary: #16213e;
    --bg-secondary: #0f3460;
    --text-primary: #f1f2f6;
    --text-secondary: #dfe6e9;
    --header-bg: #0f3460;
    --header-text: #ffffff;
    --nav-bg: #16213e;
    --nav-text: #ffffff;
    --nav-active: #e94560;
    --border-color: #303952;
    --shadow-color: rgba(0, 0, 0, 0.3);
    --accent-color: #0984e3;
    --accent-hover: #74b9ff;
    --neutral-color: #b2bec3;
    --neutral-hover: #dfe6e9;
    --success-color: #00b894;
    --success-bg: rgba(0, 184, 148, 0.2);
    --warning-color: #fdcb6e;
    --error-color: #ff7675;
    --error-bg: rgba(255, 118, 117, 0.2);
    --info-bg: rgba(9, 132, 227, 0.2);
  }

  /* Global styles */
  :global(body) {
    margin: 0;
    padding: 0;
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen,
      Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif;
    transition: background-color 0.3s, color 0.3s;
  }
  
  :global(body.dark-mode) {
    background-color: var(--bg-main);
    color: var(--text-primary);
  }
  
  main {
    background-color: var(--bg-main);
    color: var(--text-primary);
    min-height: 100vh;
  }

  .app-container {
    display: flex;
    flex-direction: column;
    min-height: 100vh;
  }

  /* Header styles */
  .app-header {
    background-color: var(--header-bg);
    color: var(--header-text);
    padding: 1rem 2rem;
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .app-header h1 {
    margin: 0;
    font-size: 1.5rem;
  }
  
  .header-right {
    display: flex;
    align-items: center;
    gap: 1rem;
  }
  
  .theme-toggle {
    background: transparent;
    border: none;
    color: var(--header-text);
    font-size: 1.2rem;
    cursor: pointer;
    padding: 0.25rem;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 50%;
    width: 2rem;
    height: 2rem;
    transition: background-color 0.3s;
  }
  
  .theme-toggle:hover {
    background-color: rgba(255, 255, 255, 0.1);
  }

  .version {
    font-size: 0.8rem;
    opacity: 0.8;
  }

  /* Navigation styles */
  .app-nav {
    background-color: var(--nav-bg);
    color: var(--nav-text);
  }

  .app-nav ul {
    display: flex;
    list-style: none;
    margin: 0;
    padding: 0;
  }

  .app-nav li {
    padding: 0;
  }

  .app-nav button {
    background: none;
    border: none;
    color: var(--nav-text);
    padding: 1rem 2rem;
    font-size: 1rem;
    cursor: pointer;
    transition: background-color 0.3s;
    width: 100%;
    text-align: center;
  }

  .app-nav li.active button {
    background-color: var(--nav-active);
  }

  .app-nav button:hover {
    background-color: var(--nav-active);
  }

  /* Content area */
  .app-content {
    flex: 1;
    padding: 2rem;
    overflow-y: auto;
  }

  /* Footer styles */
  .app-footer {
    background-color: var(--header-bg);
    color: var(--header-text);
    text-align: center;
    padding: 1rem;
    font-size: 0.8rem;
  }
</style>
