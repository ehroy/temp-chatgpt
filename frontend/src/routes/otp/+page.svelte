<script lang="ts">
	import OtpSearchForm from '../../components/OtpSearchForm.svelte';
	import OtpResultCard from '../../components/OtpResultCard.svelte';
	import StatusAlert from '../../components/StatusAlert.svelte';
	import { lookupOtp } from '$lib/api/otp';
	import type { OtpLookupResponse } from '$lib/types/otp';

	let email = '';
	let loading = false;
	let result: OtpLookupResponse | null = null;

	async function handleLookup() {
		loading = true;
		result = await lookupOtp(email);
		loading = false;
	}
</script>

<main class="shell">
	<section class="panel">
		<h1>OTP Search</h1>
		<OtpSearchForm bind:value={email} {loading} on:submit={handleLookup} />
		{#if result}
			<StatusAlert status={result.status} message={result.message} />
			<OtpResultCard {result} />
		{/if}
	</section>
</main>

<style>
	.shell { max-width: 760px; margin: 0 auto; padding: 3rem 1.2rem; }
	.panel { padding: 1.2rem; border-radius: 24px; background: rgba(8, 15, 27, 0.78); border: 1px solid rgba(148, 163, 184, 0.18); }
</style>
