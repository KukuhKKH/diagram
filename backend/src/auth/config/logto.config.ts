import { Injectable } from '@nestjs/common';
import { ConfigService } from '@nestjs/config';

@Injectable()
export class LogtoConfigService {
   private readonly appId: string;
   private readonly appSecret: string;
   private readonly endpoint: string;
   private readonly redirectUri: string;

   constructor(private readonly configService: ConfigService) {
      this.appId = this.configService.getOrThrow('LOGTO_APP_ID');
      this.appSecret = this.configService.getOrThrow('LOGTO_APP_SECRET');
      this.endpoint = this.configService.getOrThrow('LOGTO_ENDPOINT');
      this.redirectUri = this.configService.getOrThrow('LOGTO_REDIRECT_URI');
   }

   getAppId(): string {
      return this.appId;
   }

   getAppSecret(): string {
      return this.appSecret;
   }

   getEndpoint(): string {
      return this.endpoint;
   }

   getRedirectUri(): string {
      return this.redirectUri;
   }

   getOAuthConfig() {
      return {
         clientID: this.appId,
         clientSecret: this.appSecret,
         authorizationURL: `${this.endpoint}/oidc/auth`,
         tokenURL: `${this.endpoint}/oidc/token`,
         userInfoURL: `${this.endpoint}/oidc/me`,
         callbackURL: this.redirectUri,
         scope: ['openid', 'profile', 'email'],
      };
   }
}
