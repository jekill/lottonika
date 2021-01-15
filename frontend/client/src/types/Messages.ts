import { UUID } from '@/types/GenericTypes';
import { CardDto } from '@/models/CardDto';

interface CommonMessage {
  id: UUID;
  type: string;
  payload: any;
}

export interface RefreshMessage extends CommonMessage {
  type: 'refresh';
  payload: {
    card: CardDto;
  };
}

export interface RoundMessage extends CommonMessage {
  type: 'round';
  payload: {
    isWin: boolean;
    card: CardDto;
  };
}

export type CommunicationMessages = RoundMessage | RefreshMessage;
