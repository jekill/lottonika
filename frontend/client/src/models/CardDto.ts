import {UUID} from '@/types/GenericTypes';

export interface CardDto {
  id: UUID;
  number: number;
  is_win: boolean;
  is_closed: boolean;
}
