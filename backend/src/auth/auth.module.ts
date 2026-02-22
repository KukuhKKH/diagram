import { Module, MiddlewareConsumer, NestModule, Logger } from '@nestjs/common';
import { PassportModule } from '@nestjs/passport';
import { JwtModule } from '@nestjs/jwt';
import { ConfigService } from '@nestjs/config';

import { AuthController } from './auth.controller';
import { AuthService } from './auth.service';
import { LogtoStrategy } from './strategies/logto.strategy';
import { LogtoConfigService } from './config/logto.config';
import { SessionMiddleware } from './session/session.middleware';
import { ISessionStore } from './session/session-store.interface';
import { LocalSessionStore } from './session/local-session.store';
import { RedisSessionStore } from './session/redis-session.store';
import { SessionConfigService } from './session/session.config';
import { PrismaModule } from '../prisma/prisma.module';

@Module({
   imports: [
      PrismaModule,
      PassportModule.register({ defaultStrategy: 'logto' }),

      JwtModule.registerAsync({
         inject: [ConfigService],
         useFactory: (configService: ConfigService) => ({
            secret: configService.getOrThrow('JWT_SECRET'),
            signOptions: { expiresIn: '24h' },
         }),
      }),
   ],
   controllers: [AuthController],
   providers: [
      AuthService,
      LogtoStrategy,
      LogtoConfigService,
      SessionConfigService,
      {
         provide: ISessionStore,
         useFactory: (
            sessionConfig: SessionConfigService,
            configService: ConfigService,
         ): ISessionStore => {
            const logger = new Logger('SessionStoreFactory');
            const storeType = sessionConfig.getStoreType();

            logger.log(`Initializing session store: ${storeType}`);

            if (storeType === 'redis') {
               logger.warn(
                  'Redis session store selected. Make sure ioredis is installed.',
               );
               return new RedisSessionStore(configService);
            }

            return new LocalSessionStore();
         },

         inject: [SessionConfigService, ConfigService],
      },
   ],

   exports: [AuthService, LogtoConfigService, ISessionStore],
})
export class AuthModule implements NestModule {
   configure(consumer: MiddlewareConsumer) {
      consumer.apply(SessionMiddleware).forRoutes('*');
   }
}
