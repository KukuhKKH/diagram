import { Injectable, Logger } from '@nestjs/common';
import { PassportStrategy } from '@nestjs/passport';
import { Strategy, VerifyCallback } from 'passport-oauth2';
import type { StrategyOptions } from 'passport-oauth2';
import { LogtoConfigService } from '../config/logto.config';
import type { LogtoUserProfile, OAuthTokens } from '../types/session.types';
import axios from 'axios';
import { validateLogtoProfile } from '../validators/logto.validator';

@Injectable()
export class LogtoStrategy extends PassportStrategy(Strategy, 'logto') {
   private readonly logger = new Logger(LogtoStrategy.name);

   constructor(private readonly logtoConfig: LogtoConfigService) {
      const oauthConfig = logtoConfig.getOAuthConfig();
      const strategyOptions: StrategyOptions = {
         clientID: oauthConfig.clientID,
         clientSecret: oauthConfig.clientSecret,
         authorizationURL: oauthConfig.authorizationURL,
         tokenURL: oauthConfig.tokenURL,
         callbackURL: oauthConfig.callbackURL,
         state: true,
      };

      // eslint-disable-next-line @typescript-eslint/no-unsafe-call
      super(strategyOptions);
   }

   /**
    * Verify callback: called after successful token exchange.
    * Fetches user profile from Logto and returns.
    */
   async validate(
      accessToken: string,
      refreshToken: string | undefined,
      params: Record<string, unknown>,
      profile: unknown,
      done: VerifyCallback,
   ) {
      try {
         const userProfile = await this.fetchUserProfile(accessToken);

         const tokens: OAuthTokens = {
            accessToken,
            refreshToken,
            idToken: params.id_token as string | undefined,
            expiresIn: params.expires_in as number | undefined,
         };

         const cb = done as (err: Error | null, user?: unknown) => void;
         cb(null, { userProfile, tokens });
      } catch (error) {
         this.logger.error('Failed to fetch Logto user profile', error);
         const cb = done as (err: Error | null) => void;
         if (error instanceof Error) cb(error);
         else cb(new Error(String(error)));
      }
   }

   /**
    * Fetch authenticated user profile from Logto.
    * Implements the /oidc/me endpoint.
    */
   private async fetchUserProfile(
      accessToken: string,
   ): Promise<LogtoUserProfile> {
      const userInfoURL = `${this.logtoConfig.getEndpoint()}/oidc/me`;

      // eslint-disable-next-line @typescript-eslint/no-unsafe-assignment, @typescript-eslint/no-unsafe-call, @typescript-eslint/no-unsafe-member-access
      const response = await axios.get<unknown>(userInfoURL, {
         headers: {
            Authorization: `Bearer ${accessToken}`,
         },
      });

      // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access
      const validated = await validateLogtoProfile(response.data);

      return validated;
   }

   /**
    * Override authorizationParams to add custom scopes or parameters.
    * Can be extended for custom Logto configuration.
    */
   authorizationParams(options: Record<string, unknown>) {
      return {
         ...options,
      };
   }
}
