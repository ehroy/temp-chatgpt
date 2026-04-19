<script lang="ts">
	import type { OtpLookupResponse } from '$lib/types/otp';
	import { formatTime } from '$lib/utils/formatTime';

	export let result: OtpLookupResponse | null;

	const statusLabel: Record<string, string> = {
		found: 'Found',
		not_found: 'Not found',
		expired: 'Expired',
		invalid_email: 'Invalid email',
		error: 'Error'
	};

	const statusTone: Record<string, string> = {
		found: 'success',
		not_found: 'neutral',
		expired: 'warning',
		invalid_email: 'danger',
		error: 'danger'
	};
</script>

{#if result}
<section class="card">
	<div class="header">
		<div>
			<p class="eyebrow">Lookup result</p>
			<h3>{result.subject ?? 'Email matched'}</h3>
			<p class="summary">{result.message}</p>
		</div>
		<span class={`badge ${statusTone[result.status] ?? 'neutral'}`}>
			{statusLabel[result.status] ?? result.status}
		</span>
	</div>

	<div class="meta-grid">
		<div class="meta">
			<span class="label">Email</span>
			<strong>{result.email}</strong>
		</div>
		{#if result.otp}
			<div class="meta accent">
				<span class="label">OTP</span>
				<strong class="otp">{result.otp}</strong>
			</div>
		{/if}
		<div class="meta">
			<span class="label">Sender</span>
			<strong>{result.sender ?? '-'}</strong>
		</div>
		<div class="meta">
			<span class="label">Folder</span>
			<strong>{result.folder ?? '-'}</strong>
		</div>
		<div class="meta">
			<span class="label">Diterima</span>
			<strong>{formatTime(result.receivedAt)}</strong>
		</div>
		<div class="meta">
			<span class="label">Kadaluarsa</span>
			<strong>{formatTime(result.expiresAt)}</strong>
		</div>
	</div>


</section>
{/if}

<style>
	.card {
		margin-top: 1rem;
		padding: 1.1rem;
		border-radius: 28px;
		background: linear-gradient(180deg, rgba(8, 15, 27, 0.92), rgba(8, 15, 27, 0.76));
		border: 1px solid rgba(148, 163, 184, 0.14);
		box-shadow: 0 28px 70px rgba(2, 6, 23, 0.28);
	}
	.header {
		display: flex;
		justify-content: space-between;
		gap: 1rem;
		align-items: start;
		padding-bottom: 1rem;
		border-bottom: 1px solid rgba(148, 163, 184, 0.12);
	}
	.eyebrow {
		margin: 0 0 0.3rem;
		text-transform: uppercase;
		letter-spacing: 0.16em;
		font-size: 0.72rem;
		color: #93c5fd;
		font-weight: 700;
	}
	h3 {
		margin: 0;
		font-size: 1.1rem;
		letter-spacing: -0.02em;
	}
	.summary {
		margin: 0.45rem 0 0;
		color: #cbd5e1;
	}
	.badge {
		padding: 0.45rem 0.75rem;
		border-radius: 999px;
		font-size: 0.78rem;
		font-weight: 700;
		border: 1px solid transparent;
		white-space: nowrap;
	}
	.badge.success { background: rgba(52, 211, 153, 0.14); border-color: rgba(52, 211, 153, 0.26); color: #86efac; }
	.badge.warning { background: rgba(251, 191, 36, 0.14); border-color: rgba(251, 191, 36, 0.26); color: #fde68a; }
	.badge.danger { background: rgba(248, 113, 113, 0.14); border-color: rgba(248, 113, 113, 0.26); color: #fca5a5; }
	.badge.neutral { background: rgba(148, 163, 184, 0.12); border-color: rgba(148, 163, 184, 0.2); color: #e2e8f0; }
	.meta-grid {
		display: grid;
		grid-template-columns: repeat(2, minmax(0, 1fr));
		gap: 0.85rem;
		margin-top: 1rem;
	}
	.meta {
		padding: 0.95rem 1rem;
		border-radius: 18px;
		background: rgba(15, 23, 42, 0.55);
		border: 1px solid rgba(148, 163, 184, 0.12);
	}
	.meta.accent {
		background: linear-gradient(135deg, rgba(37, 99, 235, 0.18), rgba(34, 197, 94, 0.12));
	}
	.meta .label {
		display: block;
		margin-bottom: 0.3rem;
		font-size: 0.72rem;
		text-transform: uppercase;
		letter-spacing: 0.12em;
		color: #94a3b8;
	}
	.meta strong {
		font-weight: 600;
		line-height: 1.45;
		word-break: break-word;
	}
	.preview {
		width: 100%;
		min-height: 640px;
		border: 0;
		background: #fff;
	}
	.viewer {
		margin-top: 1rem;
		border-radius: 20px;
		overflow: hidden;
		background: #fff;
		border: 1px solid rgba(148, 163, 184, 0.18);
		box-shadow: 0 24px 60px rgba(0, 0, 0, 0.18);
	}
	.viewer-bar {
		display: flex;
		align-items: center;
		gap: 0.45rem;
		padding: 0.75rem 1rem;
		background: linear-gradient(180deg, #f8fafc, #eef2f7);
		border-bottom: 1px solid rgba(148, 163, 184, 0.16);
	}
	.viewer-title {
		margin-left: 0.5rem;
		font-size: 0.8rem;
		color: #64748b;
	}
	.dot {
		width: 10px;
		height: 10px;
		border-radius: 999px;
		display: inline-block;
	}
	.red { background: #ef4444; }
	.yellow { background: #f59e0b; }
	.green { background: #10b981; }
	.viewer .preview {
		display: block;
	}
	@media (max-width: 640px) {
		.preview { min-height: 560px; }
		.meta-grid { grid-template-columns: 1fr; }
		.header { flex-direction: column; }
	}
	.otp {
		font-size: 1.35rem;
		font-weight: 700;
		letter-spacing: 0.12em;
	}
</style>
