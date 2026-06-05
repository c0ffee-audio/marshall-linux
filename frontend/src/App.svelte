<script>
  import { Connect, Disconnect, GetInfo, SetANC, SetEQ } from "../wailsjs/go/main/App";

  let target = "MOTIF II A.N.C.";
  let connecting = false;
  let connected = false;
  let error = "";

  let model = "";
  let firmware = "";
  let battery = 0;
  let ancMode = "";
  let activeEQ = "";

  const ancModes = [
    { id: "off",          label: "Off" },
    { id: "anc",          label: "ANC" },
    { id: "transparency", label: "Ambient" },
  ];

  const eqPresets = [
    { id: "flat",              label: "Flat" },
    { id: "rock",              label: "Rock" },
    { id: "metal",             label: "Metal" },
    { id: "pop",               label: "Pop" },
    { id: "hiphop",            label: "Hip-Hop" },
    { id: "electronic",        label: "Electronic" },
    { id: "jazz",              label: "Jazz" },
    { id: "bass-boost",        label: "Bass" },
    { id: "mid-boost",         label: "Mid" },
    { id: "treble-boost",      label: "Treble" },
    { id: "loud-push-workout", label: "Workout" },
  ];

  async function connect() {
    connecting = true;
    error = "";
    try {
      await Connect(target);
      connected = true;
      await refresh();
    } catch (e) {
      error = e.toString();
    } finally {
      connecting = false;
    }
  }

  async function disconnect() {
    await Disconnect();
    connected = false;
    model = firmware = ancMode = activeEQ = "";
    battery = 0;
  }

  async function refresh() {
    try {
      const info = await GetInfo();
      model    = info.model;
      firmware = info.firmware;
      battery  = info.battery;
      ancMode  = info.anc;
    } catch (e) {
      error = e.toString();
    }
  }

  async function setANC(mode) {
    error = "";
    try {
      await SetANC(mode);
      ancMode = mode;
    } catch (e) {
      error = e.toString();
    }
  }

  async function setEQ(preset) {
    error = "";
    try {
      await SetEQ(preset);
      activeEQ = preset;
    } catch (e) {
      error = e.toString();
    }
  }
</script>

<main>

  <!-- Header -->
  <header>
    <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
      <path d="M9 18V5l12-2v13"/>
      <circle cx="6" cy="18" r="3"/>
      <circle cx="18" cy="16" r="3"/>
    </svg>
    <span>Marshall Linux</span>
    {#if connected}
      <div class="status-dot"></div>
    {/if}
  </header>

  {#if !connected}

    <!-- Connect -->
    <section class="connect-section">
      <p class="hint">Nom de l'appareil ou adresse MAC</p>
      <div class="input-row">
        <input
          bind:value={target}
          placeholder="MOTIF II A.N.C."
          on:keydown={(e) => e.key === "Enter" && connect()}
        />
        <button class="btn-connect" on:click={connect} disabled={connecting}>
          {#if connecting}
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="spin">
              <path d="M21 12a9 9 0 1 1-6.219-8.56"/>
            </svg>
          {:else}
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M5 12h14M12 5l7 7-7 7"/>
            </svg>
          {/if}
        </button>
      </div>
      {#if error}<p class="error">{error}</p>{/if}
    </section>

  {:else}

    <!-- Device info -->
    <section class="device-card">
      <div class="device-left">
        <svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.2">
          <path d="M3 18v-6a9 9 0 0 1 18 0v6"/>
          <path d="M21 19a2 2 0 0 1-2 2h-1a2 2 0 0 1-2-2v-3a2 2 0 0 1 2-2h3zM3 19a2 2 0 0 0 2 2h1a2 2 0 0 0 2-2v-3a2 2 0 0 0-2-2H3z"/>
        </svg>
        <div>
          <div class="device-name">{model}</div>
          <div class="device-fw">{firmware}</div>
        </div>
      </div>
      <div class="battery">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8">
          <rect x="2" y="7" width="18" height="11" rx="2"/>
          <path d="M22 11v3"/>
          <rect x="4" y="9" width="{Math.round(battery / 100 * 14)}" height="7" rx="1" fill="currentColor" stroke="none"
            class:low={battery <= 25}
            class:med={battery > 25 && battery <= 60}
            class:full={battery > 60}
          />
        </svg>
        <span class:low={battery <= 25} class:med={battery > 25 && battery <= 60}>{battery}%</span>
      </div>
    </section>

    <!-- ANC -->
    <section class="block">
      <div class="block-header">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8">
          <path d="M11 5 6 9H2v6h4l5 4zM15.54 8.46a5 5 0 0 1 0 7.07M19.07 4.93a10 10 0 0 1 0 14.14"/>
        </svg>
        Noise Control
      </div>
      <div class="anc-row">
        {#each ancModes as m}
          <button class="seg-btn" class:active={ancMode === m.id} on:click={() => setANC(m.id)}>
            {m.label}
          </button>
        {/each}
      </div>
    </section>

    <!-- EQ -->
    <section class="block">
      <div class="block-header">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8">
          <line x1="4" y1="21" x2="4" y2="14"/><line x1="4" y1="10" x2="4" y2="3"/>
          <line x1="12" y1="21" x2="12" y2="12"/><line x1="12" y1="8" x2="12" y2="3"/>
          <line x1="20" y1="21" x2="20" y2="16"/><line x1="20" y1="12" x2="20" y2="3"/>
          <line x1="1" y1="14" x2="7" y2="14"/><line x1="9" y1="8" x2="15" y2="8"/><line x1="17" y1="16" x2="23" y2="16"/>
        </svg>
        Equalizer
      </div>
      <div class="eq-grid">
        {#each eqPresets as p}
          <button class="eq-btn" class:active={activeEQ === p.id} on:click={() => setEQ(p.id)}>
            {p.label}
          </button>
        {/each}
      </div>
    </section>

    <!-- Footer -->
    <footer>
      <button class="icon-btn" title="Actualiser" on:click={refresh}>
        <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8">
          <path d="M23 4v6h-6"/><path d="M1 20v-6h6"/>
          <path d="M3.51 9a9 9 0 0 1 14.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0 0 20.49 15"/>
        </svg>
      </button>
      <button class="icon-btn disconnect" title="Déconnecter" on:click={disconnect}>
        <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8">
          <path d="M18.36 6.64A9 9 0 1 1 5.64 19.36"/><path d="M12 2v10"/>
        </svg>
      </button>
    </footer>

    {#if error}<p class="error">{error}</p>{/if}

  {/if}
</main>

<style>
  :global(*, *::before, *::after) { box-sizing: border-box; margin: 0; padding: 0; }
  :global(body) {
    background: #111;
    color: #e0e0e0;
    font-family: 'Inter', system-ui, sans-serif;
    font-size: 13px;
    -webkit-font-smoothing: antialiased;
    user-select: none;
  }

  main {
    padding: 24px 20px;
    display: flex;
    flex-direction: column;
    gap: 24px;
    min-height: 100vh;
  }

  /* Header */
  header {
    display: flex;
    align-items: center;
    gap: 10px;
    color: #fff;
    font-size: 14px;
    font-weight: 600;
    letter-spacing: 0.01em;
  }
  header svg { opacity: 0.7; }
  .status-dot {
    width: 6px; height: 6px;
    background: #e8c84a;
    border-radius: 50%;
    margin-left: auto;
    animation: pulse 3s infinite;
  }
  @keyframes pulse { 0%,100%{opacity:1} 50%{opacity:.3} }

  /* Connect */
  .connect-section { display: flex; flex-direction: column; gap: 10px; }
  .hint { color: #444; font-size: 11px; text-transform: uppercase; letter-spacing: 0.06em; }

  .input-row { display: flex; gap: 8px; }

  input {
    flex: 1;
    background: #1a1a1a;
    border: 1px solid #252525;
    border-radius: 6px;
    padding: 9px 12px;
    color: #e0e0e0;
    font-size: 13px;
    outline: none;
    transition: border-color 0.15s;
  }
  input:focus { border-color: #e8c84a; }

  .btn-connect {
    background: #e8c84a;
    border: none;
    border-radius: 6px;
    width: 38px;
    display: flex; align-items: center; justify-content: center;
    cursor: pointer;
    color: #111;
    flex-shrink: 0;
    transition: opacity 0.15s;
  }
  .btn-connect:disabled { opacity: 0.4; cursor: not-allowed; }
  .btn-connect:hover:not(:disabled) { opacity: 0.85; }

  .spin { animation: spin 0.8s linear infinite; }
  @keyframes spin { to { transform: rotate(360deg); } }

  /* Device card */
  .device-card {
    background: #161616;
    border: 1px solid #1e1e1e;
    border-radius: 10px;
    padding: 14px 16px;
    display: flex;
    align-items: center;
    justify-content: space-between;
  }
  .device-left { display: flex; align-items: center; gap: 12px; color: #555; }
  .device-left svg { flex-shrink: 0; }
  .device-name { color: #e0e0e0; font-weight: 600; font-size: 13px; margin-bottom: 2px; }
  .device-fw { color: #3a3a3a; font-size: 10px; font-family: monospace; }

  .battery { display: flex; align-items: center; gap: 5px; font-size: 12px; color: #4ade80; }
  .battery .low, .battery.low { color: #f87171; }
  .battery .med, .battery.med { color: #facc15; }

  /* Blocks */
  .block { display: flex; flex-direction: column; gap: 12px; }
  .block-header {
    display: flex; align-items: center; gap: 7px;
    color: #555;
    font-size: 10px;
    text-transform: uppercase;
    letter-spacing: 0.08em;
    font-weight: 600;
  }

  /* ANC segmented */
  .anc-row {
    display: flex;
    background: #161616;
    border: 1px solid #1e1e1e;
    border-radius: 8px;
    overflow: hidden;
  }
  .seg-btn {
    flex: 1;
    background: transparent;
    border: none;
    padding: 10px 0;
    color: #666;
    font-size: 12px;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.15s;
    letter-spacing: 0.02em;
  }
  .seg-btn:hover { color: #aaa; }
  .seg-btn.active { background: #1e1e1e; color: #e8c84a; }
  .seg-btn + .seg-btn { border-left: 1px solid #1e1e1e; }

  /* EQ grid */
  .eq-grid {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 5px;
  }
  .eq-btn {
    background: #161616;
    border: 1px solid #1e1e1e;
    border-radius: 6px;
    padding: 7px 4px;
    color: #666;
    font-size: 11px;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.15s;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }
  .eq-btn:hover { color: #bbb; border-color: #2a2a2a; }
  .eq-btn.active { color: #e8c84a; border-color: #e8c84a30; background: #e8c84a08; }

  /* Footer */
  footer { display: flex; gap: 8px; justify-content: flex-end; margin-top: auto; }

  .icon-btn {
    background: #161616;
    border: 1px solid #1e1e1e;
    border-radius: 6px;
    width: 32px; height: 32px;
    display: flex; align-items: center; justify-content: center;
    cursor: pointer;
    color: #666;
    transition: all 0.15s;
  }
  .icon-btn:hover { color: #888; border-color: #2a2a2a; }
  .icon-btn.disconnect:hover { color: #f87171; border-color: #f8717130; }

  .error {
    color: #f87171;
    font-size: 11px;
    padding: 8px 10px;
    background: #1e1515;
    border-radius: 6px;
    border: 1px solid #f8717120;
  }
</style>
