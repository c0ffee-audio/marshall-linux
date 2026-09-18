<script>
  import willenImg from "../assets/images/devices/willen.webp";
  import embertonImg from "../assets/images/devices/emberton.webp";
  import kilburnImg from "../assets/images/devices/kilburn.webp";
  import middletonImg from "../assets/images/devices/middleton.webp";
  import stanmoreImg from "../assets/images/devices/stanmore.webp";
  import woburnImg from "../assets/images/devices/woburn.webp";
  import actonImg from "../assets/images/devices/acton.webp";
  import tuftonImg from "../assets/images/devices/tufton.webp";
  import stockwellImg from "../assets/images/devices/stockwell.webp";
  import majorImg from "../assets/images/devices/major.webp";
  import monitorImg from "../assets/images/devices/monitor.webp";
  import motifImg from "../assets/images/devices/motif.webp";

  export let name = "";
  export let size = 16;
  export let strokeWidth = 1.5;

  // Vraies photos produit Marshall (marshall.com), une par gamme reconnue.
  const photosByModel = [
    [/willen/, willenImg],
    [/emberton/, embertonImg],
    [/kilburn/, kilburnImg],
    [/middleton/, middletonImg],
    [/stanmore/, stanmoreImg],
    [/woburn/, woburnImg],
    [/acton/, actonImg],
    [/tufton/, tuftonImg],
    [/stockwell/, stockwellImg],
    [/major/, majorImg],
    [/monitor/, monitorImg],
    [/motif/, motifImg],
  ];

  // Repli vectoriel par famille de forme, pour un modèle non reconnu
  // (ex. Uxbridge, discontinué) ou si la photo ne charge pas.
  function familyOf(n) {
    const s = (n || "").toLowerCase();
    if (/willen/.test(s)) return "puck";
    if (/emberton/.test(s)) return "capsule";
    if (/kilburn|middleton/.test(s)) return "portable";
    if (/stanmore|woburn|acton/.test(s)) return "homespeaker";
    if (/uxbridge|tufton|stockwell/.test(s)) return "smart";
    return "headphone";
  }

  function photoOf(n) {
    const s = (n || "").toLowerCase();
    for (const [re, img] of photosByModel) if (re.test(s)) return img;
    return null;
  }

  let imgFailed = false;
  $: {
    void name;
    imgFailed = false;
  }
  $: photo = photoOf(name);
  $: family = familyOf(name);
</script>

{#if photo && !imgFailed}
  <img
    src={photo}
    alt={name}
    width={size}
    height={size}
    style="width:{size}px; height:{size}px; object-fit:contain;"
    on:error={() => (imgFailed = true)}
  />
{:else}
  <svg
    width={size}
    height={size}
    viewBox="0 0 24 24"
    fill="none"
    stroke="currentColor"
    stroke-width={strokeWidth}
    stroke-linecap="round"
    stroke-linejoin="round"
    style="flex-shrink: 0;"
  >
    {#if family === "puck"}
      <!-- Willen: enceinte "puck" plate avec anse -->
      <path d="M9 8V6.5a3 3 0 0 1 6 0V8" />
      <rect x="3" y="8" width="18" height="9" rx="4.5" />
      <circle cx="8.5" cy="12.5" r="0.55" fill="currentColor" stroke="none" />
      <circle cx="12" cy="12.5" r="0.55" fill="currentColor" stroke="none" />
      <circle cx="15.5" cy="12.5" r="0.55" fill="currentColor" stroke="none" />

    {:else if family === "capsule"}
      <!-- Emberton: enceinte capsule verticale -->
      <rect x="7" y="2" width="10" height="20" rx="5" />
      <path d="M9 5.5h6" />
      <circle cx="10.5" cy="12" r="0.55" fill="currentColor" stroke="none" />
      <circle cx="13.5" cy="12" r="0.55" fill="currentColor" stroke="none" />
      <circle cx="10.5" cy="15" r="0.55" fill="currentColor" stroke="none" />
      <circle cx="13.5" cy="15" r="0.55" fill="currentColor" stroke="none" />

    {:else if family === "portable"}
      <!-- Kilburn / Middleton: enceinte portable avec anse -->
      <path d="M7 7V5a2 2 0 0 1 2-2h6a2 2 0 0 1 2 2v2" />
      <rect x="3" y="7" width="18" height="12" rx="2" />
      <circle cx="7.5" cy="11" r="1" />
      <line x1="12" y1="9.5" x2="12" y2="16.5" />
      <line x1="15" y1="9.5" x2="15" y2="16.5" />
      <line x1="18" y1="9.5" x2="18" y2="16.5" />

    {:else if family === "homespeaker"}
      <!-- Stanmore / Woburn / Acton: enceinte de salon rétro -->
      <rect x="2" y="4" width="20" height="16" rx="2" />
      <circle cx="7.5" cy="8" r="1.2" />
      <circle cx="16.5" cy="8" r="1.2" />
      <rect x="5.5" y="11" width="13" height="6.5" rx="1.2" />

    {:else if family === "smart"}
      <!-- Uxbridge / Tufton / Stockwell: enceinte connectée -->
      <path d="M9.2 20h5.6l1.6-12.5A2 2 0 0 0 14.4 5H9.6a2 2 0 0 0-2 2.5z" />
      <line x1="8.6" y1="11" x2="15.4" y2="11" />
      <circle cx="12" cy="3.2" r="0.6" fill="currentColor" stroke="none" />

    {:else}
      <!-- Casque (Major / Monitor / Motif / défaut) -->
      <path d="M3 18v-6a9 9 0 0 1 18 0v6" />
      <path d="M21 19a2 2 0 0 1-2 2h-1a2 2 0 0 1-2-2v-3a2 2 0 0 1 2-2h3zM3 19a2 2 0 0 0 2 2h1a2 2 0 0 0 2-2v-3a2 2 0 0 0-2-2H3z" />
    {/if}
  </svg>
{/if}
