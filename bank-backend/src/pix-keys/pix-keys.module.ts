import { Module } from '@nestjs/common';
import { PixKeysService } from './pix-keys.service';
import { PixKeysController } from './pix-keys.controller';
import { TypeOrmModule } from '@nestjs/typeorm';
import { PixKey } from './entities/pix-key.entity';
import { BankAccount } from 'src/bank-accounts/entities/bank-account.entity';
import { ClientsModule, Transport } from '@nestjs/microservices';
import { join } from 'path';
import { ConfigService } from '@nestjs/config';

@Module({
  imports: [
    TypeOrmModule.forFeature([PixKey, BankAccount]),
    ClientsModule.registerAsync([
      {
        name: 'PIX_PACKAGE',
        useFactory: (configService: ConfigService) => ({
          transport: Transport.GRPC,
          options: {
            url: configService.get('GRPC_URL'),
            package: 'derso.com.imersao.fullcycle.codeflix',
            protoPath: join(__dirname, 'proto', 'pixkey.proto'),
          },
        }),
        inject: [ConfigService],
      },
    ]),
  ],
  controllers: [PixKeysController],
  providers: [PixKeysService],
})
export class PixKeysModule {}
