import type { OtpLookupResponse } from '$lib/types/otp';

const baseUrl = import.meta.env.VITE_API_BASE_URL ?? 'http://localhost:9001';
const token = import.meta.env.VITE_API_TOKEN ?? 'dev-token';

export async function lookupOtp(email: string): Promise<OtpLookupResponse> {
	const response = await fetch(`${baseUrl}/api/otp/lookup`, {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json',
			Authorization: `Bearer ${token}`
		},
		body: JSON.stringify({ email })
	});

	const payload = (await response.json().catch(() => ({}))) as OtpLookupResponse;
	if (!response.ok) {
		return {
			status: payload.status ?? 'error',
			message: payload.message ?? 'request gagal',
			email
		};
	}

	return payload;
}
