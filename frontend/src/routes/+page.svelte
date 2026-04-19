<script>
	import OtpSearchForm from '../components/OtpSearchForm.svelte';
	import OtpResultCard from '../components/OtpResultCard.svelte';
	import StatusAlert from '../components/StatusAlert.svelte';
	import { lookupOtp } from '$lib/api/otp.js';

	let email = '';
	let loading = false;
	let result = null;

	async function handleLookup() {
		loading = true;
		result = await lookupOtp(email);
		loading = false;
	}
</script>

<main class="shell">
	<div class="background background-a"></div>
	<div class="background background-b"></div>

	<section class="hero">
		<div class="eyebrow-row">
			<span class="eyebrow">OTP Reader</span>
			<span class="pill">Yahoo IMAP</span>
		</div>
		<h1>Lookup OTP, tanpa clutter.</h1>
		<p class="lede">Cari email yang sah, lihat hasilnya dengan tampilan yang bersih, cepat, dan modern.</p>
		<div class="hero-metrics">
			<div>
				<strong>5 menit</strong>
				<span>maksimum usia</span>
			</div>
			<div>
				<strong>1 email</strong>
				<span>terbaru per lookup</span>
			</div>
			<div>
				<strong>Whitelist</strong>
				<span>folder yang diizinkan</span>
			</div>
		</div>
	</section>

	<section class="layout">
		<div class="panel search-panel">
			<div class="panel-head">
				<p class="panel-kicker">Lookup</p>
				<h2>Masukkan email</h2>
			</div>
			<OtpSearchForm bind:value={email} {loading} on:submit={handleLookup} />
			<p class="hint">Contoh: user@yahoo.com</p>
		</div>

		<div class="panel result-panel">
			<div class="panel-head">
				<p class="panel-kicker">Result</p>
				<h2>Preview email</h2>
			</div>
			{#if loading && !result}
				<div class="skeleton-card">
					<div class="skeleton-line w-40"></div>
					<div class="skeleton-line w-70"></div>
					<div class="skeleton-grid">
						<div class="skeleton-block"></div>
						<div class="skeleton-block"></div>
						<div class="skeleton-block"></div>
						<div class="skeleton-block"></div>
					</div>
				</div>
			{:else if result}
				<StatusAlert status={result.status} message={result.message} />
				<OtpResultCard {result} />
			{:else}
				<div class="empty-state">
					<p class="empty-title">Belum ada hasil</p>
					<p>Lookup akan menampilkan email terbaru yang cocok dan valid.</p>
				</div>
			{/if}
		</div>
	</section>
</main>

<style>
	.shell {
		position: relative;
		max-width: 1200px;
		margin: 0 auto;
		padding: 4rem 1.2rem 5rem;
		overflow: hidden;
	}
	.background {
		position: absolute;
		inset: auto;
		border-radius: 999px;
		filter: blur(18px);
		opacity: 0.45;
		pointer-events: none;
	}
	.background-a {
		top: -120px;
		right: -120px;
		width: 360px;
		height: 360px;
		background: radial-gradient(circle, rgba(96, 165, 250, 0.22), rgba(96, 165, 250, 0));
	}
	.background-b {
		bottom: 0;
		left: -120px;
		width: 400px;
		height: 400px;
		background: radial-gradient(circle, rgba(52, 211, 153, 0.14), rgba(52, 211, 153, 0));
	}
	.hero {
		position: relative;
		max-width: 760px;
		margin-bottom: 1.35rem;
	}
	.eyebrow-row {
		display: flex;
		align-items: center;
		gap: 0.75rem;
		flex-wrap: wrap;
		margin-bottom: 1rem;
	}
	.eyebrow {
		text-transform: uppercase;
		letter-spacing: .2em;
		color: #93c5fd;
		font-size: .72rem;
		font-weight: 700;
	}
	.pill {
		padding: 0.38rem 0.7rem;
		border-radius: 999px;
		background: rgba(148, 163, 184, 0.08);
		border: 1px solid rgba(148, 163, 184, 0.14);
		color: #cbd5e1;
		font-size: 0.78rem;
	}
	h1 {
		margin: 0.15rem 0 0.65rem;
		font-size: clamp(2.5rem, 6vw, 4.5rem);
		line-height: 0.94;
		letter-spacing: -0.055em;
		max-width: 11ch;
	}
	.lede {
		max-width: 64ch;
		color: rgba(226, 232, 240, 0.78);
		font-size: 1.03rem;
	}
	.hero-metrics {
		display: grid;
		grid-template-columns: repeat(3, minmax(0, 1fr));
		gap: 0.85rem;
		margin-top: 1.5rem;
	}
	.hero-metrics div {
		padding: 0.95rem 1rem;
		border-radius: 18px;
		background: rgba(8, 15, 27, 0.48);
		border: 1px solid rgba(148, 163, 184, 0.12);
		backdrop-filter: blur(10px);
	}
	.hero-metrics strong {
		display: block;
		font-size: 1rem;
		margin-bottom: 0.2rem;
	}
	.hero-metrics span {
		color: #94a3b8;
		font-size: 0.88rem;
	}
	.layout {
		display: grid;
		grid-template-columns: minmax(0, 380px) minmax(0, 1fr);
		gap: 1rem;
		align-items: start;
	}
	.panel {
		padding: 1.15rem;
		border-radius: 24px;
		background: rgba(7, 12, 21, 0.62);
		border: 1px solid rgba(148, 163, 184, 0.1);
		backdrop-filter: blur(14px);
		box-shadow: 0 20px 48px rgba(2, 6, 23, 0.22);
	}
	.panel-head {
		margin-bottom: 1rem;
	}
	.panel-kicker {
		margin: 0 0 0.35rem;
		color: #93c5fd;
		font-size: 0.76rem;
		text-transform: uppercase;
		letter-spacing: 0.16em;
		font-weight: 700;
	}
	h2 {
		margin: 0;
		font-size: 1.2rem;
		letter-spacing: -0.02em;
	}
	.search-panel {
		position: sticky;
		top: 1rem;
	}
	.hint {
		margin: 0.9rem 0 0;
		font-size: 0.88rem;
		color: #94a3b8;
	}
	.empty-state {
		padding: 1.3rem;
		border-radius: 20px;
		border: 1px dashed rgba(148, 163, 184, 0.18);
		background: rgba(15, 23, 42, 0.3);
	}
	.empty-title {
		margin: 0 0 0.35rem;
		font-weight: 700;
	}
	.skeleton-card {
		display: grid;
		gap: 0.9rem;
		padding: 1rem;
		border-radius: 20px;
		border: 1px solid rgba(148, 163, 184, 0.12);
		background: rgba(15, 23, 42, 0.28);
	}
	.skeleton-line,
	.skeleton-block {
		position: relative;
		overflow: hidden;
		border-radius: 999px;
		background: rgba(148, 163, 184, 0.1);
	}
	.skeleton-line {
		height: 12px;
	}
	.skeleton-grid {
		display: grid;
		grid-template-columns: repeat(2, minmax(0, 1fr));
		gap: 0.75rem;
	}
	.skeleton-block {
		height: 72px;
		border-radius: 18px;
	}
	.skeleton-line::after,
	.skeleton-block::after {
		content: '';
		position: absolute;
		inset: 0;
		transform: translateX(-100%);
		background: linear-gradient(90deg, transparent, rgba(255,255,255,0.12), transparent);
		animation: shimmer 1.3s infinite;
	}
	.w-40 { width: 40%; }
	.w-70 { width: 70%; }
	@keyframes shimmer {
		100% { transform: translateX(100%); }
	}
	@media (max-width: 920px) {
		.layout {
			grid-template-columns: 1fr;
		}
		.search-panel {
			position: static;
		}
		.hero-metrics {
			grid-template-columns: 1fr;
		}
		h1 {
			max-width: none;
		}
		.skeleton-grid {
			grid-template-columns: 1fr;
		}
	}
</style>
