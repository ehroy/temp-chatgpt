export type OtpStatus = 'found' | 'not_found' | 'expired' | 'invalid_email' | 'error';

export type OtpLookupResponse = {
	status: OtpStatus | string;
	message: string;
	email: string;
	otp?: string;
	subject?: string;
	sender?: string;
	folder?: string;
	html?: string;
	receivedAt?: string;
	expiresAt?: string;
};
