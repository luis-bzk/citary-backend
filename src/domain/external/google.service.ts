import { OAuth2Client } from 'google-auth-library';

import { EnvConfig } from '../../config';
import { GoogleUserData } from '../interfaces';

export const oAuth2Client = new OAuth2Client({
  clientId: EnvConfig().GOOGLE_CLIENT_ID!,
  clientSecret: EnvConfig().GOOGLE_CLIENT_SECRET!,
  redirectUri: EnvConfig().GOOGLE_REDIRECT_URI!,
});

export function getGoogleAuthURL() {
  const scopes = [
    'https://www.googleapis.com/auth/userinfo.profile',
    'https://www.googleapis.com/auth/userinfo.email',
  ];

  const url = oAuth2Client.generateAuthUrl({
    access_type: 'offline',
    prompt: 'consent',
    scope: scopes,
  });

  return url;
}

export async function getGoogleUser(code: string): Promise<GoogleUserData> {
  const { tokens } = await oAuth2Client.getToken(code);
  oAuth2Client.setCredentials(tokens);

  const res = await oAuth2Client.request<GoogleUserData>({
    url: 'https://www.googleapis.com/oauth2/v2/userinfo',
  });

  return res.data;
}
