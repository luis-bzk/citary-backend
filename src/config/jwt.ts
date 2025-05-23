import jwt from 'jsonwebtoken';
import { EnvConfig } from './envs';

export class JwtAdapter {
  static async generateToken(
    payload: Object,
    duration: string = '2h',
  ): Promise<string | null> {
    return new Promise((resolve) => {
      jwt.sign(
        payload,
        EnvConfig().JWT_SEED,
        { expiresIn: duration },
        (err, token) => {
          if (err) return resolve(null);
          resolve(token!);
        },
      );
    });
  }

  static validateToken<T>(token: string): Promise<T | null> {
    return new Promise((resolve) => {
      jwt.verify(token, EnvConfig().JWT_SEED, (err, decoded) => {
        if (err) return resolve(null);

        resolve(decoded as T);
      });
    });
  }
}
