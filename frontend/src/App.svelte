<script>
  import { onMount } from "svelte";
  import { Connect, Disconnect, GetInfo, SetANC, SetEQ, ScanDevices, GetCachedDevices, GetCapabilities } from "../wailsjs/go/main/App";
  import DeviceIcon from "./lib/DeviceIcon.svelte";

  let scanning = false;
  let scannedDevices = [];
  let connectingDevice = "";
  let connected = false;
  let error = "";

  let manualTarget = "";
  let showManual = false;

  let model = "";
  let firmware = "";
  let battery = 0;
  let ancMode = "";
  let activeEQ = "";

  let caps = {
    hasANC: false, hasEQ: false, hasBattery: false,
    hasVolume: false, hasRoomPlacement: false,
    hasPartyMode: false, hasNightMode: false, hasLED: false,
  };

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

  onMount(() => scan());

  async function scan() {
    scanning = true;
    scannedDevices = [];
    error = "";
    try {
      // afficher le cache BlueZ immédiatement
      const cached = await GetCachedDevices();
      if (cached.length > 0) {
        scannedDevices = cached;
      }
      // scan BLE en arrière-plan, mettre à jour si nouveaux appareils
      ScanDevices().then(fresh => {
        scannedDevices = fresh;
        if (fresh.length === 0) {
          error = "No devices found. Make sure your device is powered on.";
        }
      }).catch(() => {}).finally(() => { scanning = false; });
    } catch (e) {
      error = e.toString();
      scanning = false;
    }
  }

  async function connectTo(target) {
    connectingDevice = target;
    error = "";
    try {
      await Connect(target);
      connected = true;
      await refresh();
    } catch (e) {
      error = e.toString();
    } finally {
      connectingDevice = "";
    }
  }

  async function disconnect() {
    await Disconnect();
    connected = false;
    model = firmware = ancMode = activeEQ = "";
    battery = 0;
    caps = { hasANC: false, hasEQ: false, hasBattery: false, hasVolume: false, hasRoomPlacement: false, hasPartyMode: false, hasNightMode: false, hasLED: false };
    scannedDevices = [];
    scan();
  }

  async function refresh() {
    try {
      const [info, c] = await Promise.all([GetInfo(), GetCapabilities()]);
      model    = info.model;
      firmware = info.firmware;
      battery  = info.battery;
      ancMode  = info.anc;
      caps     = c;
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
    <div class="brand">
      <div class="brand-badge">
        <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.6">
          <path d="M9 18V5l12-2v13"/>
          <circle cx="6" cy="18" r="3"/>
          <circle cx="18" cy="16" r="3"/>
        </svg>
      </div>
      <span class="brand-name">Marshall&nbsp;Linux</span>
    </div>
    {#if connected}
      <div class="status-pill">
        <span class="status-dot"></span>
        Connected
      </div>
    {/if}
  </header>

  {#if !connected}

    <!-- Scan / device list -->
    <section class="connect-section">

      {#if scanning}
        <div class="scan-state">
          <svg width="26" height="26" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" class="spin">
            <path d="M21 12a9 9 0 1 1-6.219-8.56"/>
          </svg>
          <span class="scan-label">Searching for devices…</span>
        </div>

      {:else if scannedDevices.length > 0}
        <p class="hint">Select your device</p>
        <div class="device-list">
          {#each scannedDevices as d}
            <button
              class="device-item"
              on:click={() => connectTo(d.address)}
              disabled={!!connectingDevice}
            >
              <div class="device-item-icon">
                <DeviceIcon name={d.name} size={52} strokeWidth={1.4} />
              </div>
              <div class="device-item-info">
                <span class="device-item-name">{d.name}</span>
                <span class="device-item-addr">{d.address}</span>
              </div>
              {#if connectingDevice === d.address}
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="spin">
                  <path d="M21 12a9 9 0 1 1-6.219-8.56"/>
                </svg>
              {:else}
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="arrow-icon">
                  <path d="M5 12h14M12 5l7 7-7 7"/>
                </svg>
              {/if}
            </button>
          {/each}
        </div>
        <button class="btn-rescan" on:click={scan}>
          <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M23 4v6h-6"/><path d="M1 20v-6h6"/>
            <path d="M3.51 9a9 9 0 0 1 14.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0 0 20.49 15"/>
          </svg>
          Scan again
        </button>

      {:else}
        <button class="btn-scan" on:click={scan}>
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.6">
            <circle cx="12" cy="12" r="2"/>
            <path d="M16.24 7.76a6 6 0 0 1 0 8.49m-8.48-.01a6 6 0 0 1 0-8.49m11.31-2.82a10 10 0 0 1 0 14.14m-14.14 0a10 10 0 0 1 0-14.14"/>
          </svg>
          Scan for devices
        </button>
      {/if}

      <!-- Manual input -->
      <div class="manual-toggle">
        <button class="link-btn" on:click={() => showManual = !showManual}>
          {showManual ? "Hide" : "Enter manually"}
        </button>
      </div>
      {#if showManual}
        <div class="input-row">
          <input
            id="manual-target"
            bind:value={manualTarget}
            placeholder="Device name or MAC address"
            on:keydown={(e) => e.key === "Enter" && connectTo(manualTarget)}
          />
          <button class="btn-connect" on:click={() => connectTo(manualTarget)} disabled={!!connectingDevice}>
            {#if connectingDevice === manualTarget}
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
      {/if}

      {#if error}<p class="error">{error}</p>{/if}
    </section>

  {:else}

    <!-- Device hero -->
    <section class="hero">
      <div class="hero-glow"></div>
      <div class="hero-photo">
        <DeviceIcon name={model} size={148} strokeWidth={0.9} />
      </div>
      <div class="hero-name">{model}</div>
      <div class="hero-fw">FW {firmware}</div>

      {#if caps.hasBattery}
      <div class="battery-pill">
        <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8">
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
      {/if}
    </section>

    <!-- ANC -->
    {#if caps.hasANC}
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
    {/if}

    <!-- EQ -->
    {#if caps.hasEQ}
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
    {/if}

    <!-- No advanced features -->
    {#if !caps.hasANC && !caps.hasEQ}
    <p class="no-features">No advanced controls available for this device.</p>
    {/if}

    <!-- Footer -->
    <footer>
      <button class="btn-pill" on:click={refresh}>
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8">
          <path d="M23 4v6h-6"/><path d="M1 20v-6h6"/>
          <path d="M3.51 9a9 9 0 0 1 14.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0 0 20.49 15"/>
        </svg>
        Refresh
      </button>
      <button class="btn-pill disconnect" on:click={disconnect}>
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8">
          <path d="M18.36 6.64A9 9 0 1 1 5.64 19.36"/><path d="M12 2v10"/>
        </svg>
        Disconnect
      </button>
    </footer>

    {#if error}<p class="error">{error}</p>{/if}

  {/if}
</main>

<style>
  :global(*, *::before, *::after) { box-sizing: border-box; margin: 0; padding: 0; }
  :global(body) {
    background: #0d0d0d;
    color: #f2f0ec;
    font-family: 'Inter', system-ui, sans-serif;
    font-size: 13.5px;
    -webkit-font-smoothing: antialiased;
    user-select: none;
  }

  main {
    padding: 26px 24px 22px;
    display: flex;
    flex-direction: column;
    gap: 26px;
    min-height: 100vh;
  }

  /* Header */
  header {
    display: flex;
    align-items: center;
    gap: 10px;
  }
  .brand { display: flex; align-items: center; gap: 10px; }
  .brand-badge {
    width: 30px; height: 30px;
    display: flex; align-items: center; justify-content: center;
    background: #1a1a1a;
    border: 1px solid #232323;
    border-radius: 8px;
    color: #e8c84a;
  }
  .brand-name {
    font-family: 'Oswald', sans-serif;
    font-weight: 600;
    font-size: 17px;
    letter-spacing: 0.02em;
    text-transform: uppercase;
    color: #fff;
  }
  .status-pill {
    margin-left: auto;
    display: flex;
    align-items: center;
    gap: 6px;
    background: #e8c84a12;
    border: 1px solid #e8c84a2e;
    color: #e8c84a;
    font-size: 11px;
    font-weight: 600;
    letter-spacing: 0.02em;
    padding: 5px 10px;
    border-radius: 999px;
  }
  .status-dot {
    width: 6px; height: 6px;
    background: #e8c84a;
    border-radius: 50%;
    animation: pulse 3s infinite;
  }
  @keyframes pulse { 0%,100%{opacity:1} 50%{opacity:.35} }

  /* Connect section */
  .connect-section { display: flex; flex-direction: column; gap: 16px; flex: 1; }
  .hint { color: #4a4a4a; font-size: 11px; text-transform: uppercase; letter-spacing: 0.08em; font-weight: 600; }

  /* Scanning state */
  .scan-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 14px;
    padding: 64px 0;
    color: #555;
    flex: 1;
  }
  .scan-label { font-size: 13px; }

  /* Scan button */
  .btn-scan {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 10px;
    background: #161616;
    border: 1px solid #232323;
    border-radius: 12px;
    padding: 22px;
    color: #888;
    font-size: 14px;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.15s;
    width: 100%;
  }
  .btn-scan:hover { color: #e8c84a; border-color: #e8c84a40; background: #1a1a1a; }

  /* Device list */
  .device-list {
    display: flex;
    flex-direction: column;
    gap: 8px;
  }
  .device-item {
    display: flex;
    align-items: center;
    gap: 14px;
    background: #161616;
    border: 1px solid #1f1f1f;
    border-radius: 12px;
    padding: 10px 14px;
    cursor: pointer;
    transition: all 0.15s;
    width: 100%;
    text-align: left;
    color: #e0e0e0;
  }
  .device-item:hover:not(:disabled) { border-color: #e8c84a30; background: #1a1a1a; }
  .device-item:disabled { opacity: 0.5; cursor: not-allowed; }
  .device-item-icon { flex-shrink: 0; display: flex; align-items: center; justify-content: center; width: 52px; }
  .device-item-info { flex: 1; display: flex; flex-direction: column; gap: 3px; min-width: 0; }
  .device-item-name { font-size: 14px; font-weight: 600; color: #f0f0f0; }
  .device-item-addr { font-size: 10.5px; font-family: ui-monospace, 'SF Mono', monospace; color: #444; }
  .arrow-icon { color: #333; transition: color 0.15s; flex-shrink: 0; }
  .device-item:hover .arrow-icon { color: #e8c84a; }

  /* Rescan button */
  .btn-rescan {
    display: flex;
    align-items: center;
    gap: 6px;
    background: transparent;
    border: 1px solid #1e1e1e;
    border-radius: 8px;
    padding: 8px 13px;
    color: #4a4a4a;
    font-size: 11.5px;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.15s;
    align-self: flex-start;
  }
  .btn-rescan:hover { color: #999; border-color: #2a2a2a; }

  /* Manual input */
  .manual-toggle { display: flex; margin-top: auto; }
  .link-btn {
    background: none;
    border: none;
    color: #3a3a3a;
    font-size: 11.5px;
    cursor: pointer;
    padding: 0;
    transition: color 0.15s;
  }
  .link-btn:hover { color: #777; }

  .input-row { display: flex; gap: 8px; }

  input {
    flex: 1;
    background: #1a1a1a;
    border: 1px solid #262626;
    border-radius: 8px;
    padding: 11px 13px;
    color: #e0e0e0;
    font-size: 13.5px;
    outline: none;
    transition: border-color 0.15s;
  }
  input:focus { border-color: #e8c84a; }
  input:focus-visible { outline: 2px solid #e8c84a60; outline-offset: 1px; }

  .btn-connect {
    background: #e8c84a;
    border: none;
    border-radius: 8px;
    width: 42px;
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

  /* Hero */
  .hero {
    position: relative;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 6px;
    padding: 20px 16px 22px;
    background: #141414;
    border: 1px solid #1f1f1f;
    border-radius: 16px;
    overflow: hidden;
  }
  .hero-glow {
    position: absolute;
    top: -40%;
    left: 50%;
    width: 320px;
    height: 320px;
    transform: translateX(-50%);
    background: radial-gradient(circle, #e8c84a1c 0%, transparent 68%);
    pointer-events: none;
  }
  .hero-photo {
    position: relative;
    display: flex;
    align-items: center;
    justify-content: center;
    height: 148px;
    color: #3a3a3a;
  }
  .hero-name {
    position: relative;
    font-family: 'Oswald', sans-serif;
    font-weight: 600;
    font-size: 22px;
    letter-spacing: 0.01em;
    text-transform: uppercase;
    color: #fff;
    text-align: center;
    text-wrap: balance;
  }
  .hero-fw {
    position: relative;
    color: #4a4a4a;
    font-size: 10.5px;
    font-family: ui-monospace, monospace;
    letter-spacing: 0.03em;
  }
  .battery-pill {
    position: relative;
    display: flex; align-items: center; gap: 6px;
    margin-top: 10px;
    background: #1c1c1c;
    border: 1px solid #262626;
    border-radius: 999px;
    padding: 6px 12px;
    font-size: 12px;
    font-weight: 600;
    color: #4ade80;
  }
  .battery-pill .low, .battery-pill span.low { color: #f87171; }
  .battery-pill .med, .battery-pill span.med { color: #facc15; }

  /* Blocks */
  .block { display: flex; flex-direction: column; gap: 12px; }
  .block-header {
    display: flex; align-items: center; gap: 7px;
    color: #666;
    font-size: 10.5px;
    text-transform: uppercase;
    letter-spacing: 0.09em;
    font-weight: 700;
  }

  /* ANC segmented */
  .anc-row {
    display: flex;
    background: #161616;
    border: 1px solid #1f1f1f;
    border-radius: 10px;
    overflow: hidden;
  }
  .seg-btn {
    flex: 1;
    background: transparent;
    border: none;
    padding: 13px 0;
    color: #666;
    font-size: 12.5px;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.15s;
    letter-spacing: 0.02em;
  }
  .seg-btn:hover { color: #aaa; }
  .seg-btn.active { background: #e8c84a12; color: #e8c84a; }
  .seg-btn + .seg-btn { border-left: 1px solid #1f1f1f; }

  /* EQ grid */
  .eq-grid {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 8px;
  }
  .eq-btn {
    background: #161616;
    border: 1px solid #1f1f1f;
    border-radius: 9px;
    padding: 13px 6px;
    color: #777;
    font-size: 12px;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.15s;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }
  .eq-btn:hover { color: #bbb; border-color: #2a2a2a; }
  .eq-btn.active { color: #e8c84a; border-color: #e8c84a40; background: #e8c84a10; }

  /* Footer */
  footer { display: flex; gap: 10px; margin-top: auto; }

  .btn-pill {
    flex: 1;
    display: flex; align-items: center; justify-content: center; gap: 7px;
    background: #161616;
    border: 1px solid #1f1f1f;
    border-radius: 10px;
    padding: 11px 0;
    cursor: pointer;
    color: #888;
    font-size: 12.5px;
    font-weight: 600;
    transition: all 0.15s;
  }
  .btn-pill:hover { color: #ccc; border-color: #2a2a2a; }
  .btn-pill.disconnect:hover { color: #f87171; border-color: #f8717130; background: #1e1515; }

  .no-features {
    color: #3a3a3a;
    font-size: 11.5px;
    text-align: center;
    padding: 20px 0;
    margin-top: auto;
  }

  .error {
    color: #f87171;
    font-size: 11.5px;
    padding: 9px 12px;
    background: #1e1515;
    border-radius: 8px;
    border: 1px solid #f8717120;
  }
</style>
