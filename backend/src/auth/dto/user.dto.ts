export class UserResponseDto {
   id: string;
   logtoUserId: string;
   email?: string;
   name?: string;
   avatarUrl?: string;
   createdAt: Date;
   updatedAt: Date;
}

export class LoginResponseDto {
   redirectUrl: string;
   csrfToken: string;
}

export class AuthStatusDto {
   isAuthenticated: boolean;
   user?: UserResponseDto;
}
