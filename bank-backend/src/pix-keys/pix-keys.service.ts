import { Injectable } from '@nestjs/common';
import { CreatePixKeyDto } from './dto/create-pix-key.dto';
import { UpdatePixKeyDto } from './dto/update-pix-key.dto';
import { InjectRepository } from '@nestjs/typeorm';
import { Repository } from 'typeorm';
import { PixKey } from './entities/pix-key.entity';

@Injectable()
export class PixKeysService {
  constructor(
    @InjectRepository(PixKey)
    private pixKeyRepo: Repository<PixKey>,
  ) {}

  create(createPixKeyDto: CreatePixKeyDto) {
    return 'This action adds a new pixKey';
  }

  findAll() {
    return `This action returns all pixKeys`;
  }

  findOne(id: number) {
    return `This action returns a #${id} pixKey`;
  }

  update(id: number, updatePixKeyDto: UpdatePixKeyDto) {
    return `This action updates a #${id} pixKey`;
  }

  remove(id: number) {
    return `This action removes a #${id} pixKey`;
  }
}
