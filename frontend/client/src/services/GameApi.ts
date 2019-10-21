import {AxiosInstance} from 'axios';
import {CardDto} from '@/models/CardDto';

export class GameApi {
  constructor(private transport: AxiosInstance) {
  }

  public createCard(): Promise<CardDto> {
    return Promise.resolve({id: '123-123-dasd', number: 123});
    // return this.transport.post<CardDto>('/card').then((r) => r.data);
  }

  public stopGame(cardId: string): Promise<void> {
    return Promise.resolve();
    // return this.transport.delete('/card').then((r) => undefined);
  }
}
