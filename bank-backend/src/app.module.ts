import { Module } from '@nestjs/common';
import { AppController } from './app.controller';
import { AppService } from './app.service';
import { BankAccountsModule } from './bank-accounts/bank-accounts.module';
import { TypeOrmModule } from '@nestjs/typeorm';
import { BankAccount } from './bank-accounts/entities/bank-account.entity';
import { PixKeysModule } from './pix-keys/pix-keys.module';
import { PixKey } from './pix-keys/entities/pix-key.entity';
import { TransactionsModule } from './transactions/transactions.module';
import { ConfigModule, ConfigService } from '@nestjs/config';
import { Transaction } from './transactions/entities/transaction.entity';

@Module({
  imports: [
    // Deve ser o primeiro para estar disponível para injeção nos outros
    ConfigModule.forRoot({
      envFilePath: ['.env', `.bank-${process.env.BANK_CODE}.env`],
      isGlobal: true,
    }),
    // forRootAsync: é preciso aguardar o processamento da injeção do ConfigService,
    // o qual leu os arquivos .env (veja definição acima)
    TypeOrmModule.forRootAsync({
      useFactory: (configService: ConfigService) => ({
        type: configService.get('TYPEORM_CONNECTION') as any,
        host: configService.get('TYPEORM_HOST'),
        port: parseInt(configService.get('TYPEORM_PORT')),
        username: configService.get('TYPEORM_USERNAME'),
        password: configService.get('TYPEORM_PASSWORD'),
        database: configService.get('TYPEORM_DATABASE') as string,
        entities: [BankAccount, PixKey, Transaction],
        synchronize: true, // sincroniza as entidades com o DB
      }),
      inject: [ConfigService],
    }),
    BankAccountsModule,
    PixKeysModule,
    TransactionsModule,
  ],
  controllers: [AppController],
  providers: [AppService],
})
export class AppModule {}
