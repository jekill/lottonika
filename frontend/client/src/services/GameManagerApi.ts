import { AxiosInstance } from 'axios';

export class GameManagerApi {
  constructor(private transport: AxiosInstance, private wsHost: string) {
  }

  public wsStateConnect(): WebSocket {
    return new WebSocket(`${this.wsHost}${this.transport.defaults.baseURL}/state/ws`);
  }

  public startRound() {
    return this.transport.post('/actions', { action: 'startRound' });
  }
}
