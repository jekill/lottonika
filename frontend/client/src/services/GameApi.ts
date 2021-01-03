import {AxiosInstance} from 'axios';
import {CardDto} from '@/models/CardDto';
import {UUID} from '@/types/GenericTypes';

export class GameApi {
  constructor(private transport: AxiosInstance, private wsHost: string) {
  }

  public createCard(): Promise<CardDto> {
    // return Promise.resolve({id: '123-123-dasd', number: 123});
    return this.transport.post<CardDto>('/cards').then((r) => r.data);
  }

  public getCard(id: UUID): Promise<CardDto> {
    return this.transport.get<CardDto>('/cards/' + id).then((r) => r.data);
  }

  public wsConnect(cardId: UUID): WebSocket {
    return new WebSocket(`${this.wsHost}${this.transport.defaults.baseURL}/cards/${encodeURIComponent(cardId)}/ws`);
  }

  public stopGame(cardId: UUID): Promise<void> {
    return Promise.resolve();
    // return this.transport.delete('/card').then((r) => undefined);
  }
}
