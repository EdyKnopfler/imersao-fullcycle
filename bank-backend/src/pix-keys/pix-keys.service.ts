import { Inject, Injectable, OnModuleInit } from '@nestjs/common';
import { CreatePixKeyDto } from './dto/create-pix-key.dto';
import { UpdatePixKeyDto } from './dto/update-pix-key.dto';
import { InjectRepository } from '@nestjs/typeorm';
import { Repository } from 'typeorm';
import { PixKey, PixKeyKind } from './entities/pix-key.entity';
import { BankAccount } from 'src/bank-accounts/entities/bank-account.entity';
import { ClientGrpc } from '@nestjs/microservices';
import { PixKeyClientGrpc, RegisterPixKeyRpcResponse } from './proto/pix-keys.grpc';
import { lastValueFrom } from 'rxjs';

@Injectable()
export class PixKeysService implements OnModuleInit {
  // Tivemos que declarar os tipos TS correspondentes aos elementos no protofile
  private pixGrpcService: PixKeyClientGrpc;

  constructor(
    @InjectRepository(PixKey) private pixKeyRepo: Repository<PixKey>,
    @InjectRepository(BankAccount) private bankAccountRepo: Repository<BankAccount>,
    @Inject('PIX_PACKAGE') private pixGrpcPackage: ClientGrpc,
  ) {}

  onModuleInit() {
    this.pixGrpcService = this.pixGrpcPackage.getService('PixService');
  }

  async create(bankAccountId: string, createPixKeyDto: CreatePixKeyDto) {
    await this.bankAccountRepo.findOneOrFail({
      where: { id: bankAccountId },
    });

    const remotePixKey = await this.findRemotePixKey(createPixKeyDto);

    if (remotePixKey) {
      return this.createIfNotExists(bankAccountId, remotePixKey);
    }

    const createdRemotePixKey = await lastValueFrom(
      this.pixGrpcService.registerPixKey({
        ...createPixKeyDto,
        accountId: bankAccountId,
      }),
    );

    await this.pixKeyRepo.save({
      id: createdRemotePixKey.id,
      bank_account_id: bankAccountId,
      ...createPixKeyDto,
    });
  }

  private async findRemotePixKey(data: {
    key: string;
    kind: string;
  }): Promise<RegisterPixKeyRpcResponse> {
    try {
      return await lastValueFrom(this.pixGrpcService.find(data));
    } catch (e) {
      if (e.details == 'no key found') {
        return null;
      }
      console.error(e);
      throw new Error('gRPC internal error');
    }
  }

  private async createIfNotExists(
    bankAccountId: string,
    remotePixKey: RegisterPixKeyRpcResponse,
  ) {
    const hasLocalPixKey = await this.pixKeyRepo.exist({
      where: { key: remotePixKey.key },
    });

    if (hasLocalPixKey) {
      throw new PixKeyAlreadyExistsError();
    }

    return this.pixKeyRepo.save({
      id: remotePixKey.id,
      bank_account_id: bankAccountId,
      key: remotePixKey.key,
      kind: remotePixKey.kind as PixKeyKind,
    });
  }

  findAll(bankAccountId: string) {
    return this.pixKeyRepo.find({
      where: { bank_account_id: bankAccountId },
      order: { created_at: 'DESC' },
    });
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

export class PixKeyAlreadyExistsError extends Error {}