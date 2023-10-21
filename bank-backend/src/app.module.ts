import { Module } from '@nestjs/common';
import { AppController } from './app.controller';
import { AppService } from './app.service';
import { BankAccountsModule } from './bank-accounts/bank-accounts.module';
import { TypeOrmModule } from '@nestjs/typeorm';

@Module({
  imports: [
    TypeOrmModule.forRoot({
      type: 'postgres',
      host: 'db',
      database: 'bank',
      username: 'postgres',
      password: 'root',
      entities: [],
      synchronize: true,  // sincroniza as entidades com o DB
    }),
    BankAccountsModule
  ],
  controllers: [AppController],
  providers: [AppService],
})
export class AppModule {}
