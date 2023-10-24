import { Module } from '@nestjs/common';
import { AppController } from './app.controller';
import { AppService } from './app.service';
import { BankAccountsModule } from './bank-accounts/bank-accounts.module';
import { TypeOrmModule } from '@nestjs/typeorm';
import { BankAccount } from './bank-accounts/entities/bank-account.entity';
import { PixKeysModule } from './pix-keys/pix-keys.module';
import { PixKey } from './pix-keys/entities/pix-key.entity';
import { TransactionsModule } from './transactions/transactions.module';

@Module({
  imports: [
    TypeOrmModule.forRoot({
      type: 'postgres',
      host: 'db',
      database: 'bank',
      username: 'postgres',
      password: 'root',
      entities: [BankAccount, PixKey],
      synchronize: true, // sincroniza as entidades com o DB
    }),
    BankAccountsModule,
    PixKeysModule,
    TransactionsModule,
  ],
  controllers: [AppController],
  providers: [AppService],
})
export class AppModule {}
