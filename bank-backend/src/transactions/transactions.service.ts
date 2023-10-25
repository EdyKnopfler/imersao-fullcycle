import { Inject, Injectable } from '@nestjs/common';
import { CreateTransactionDto } from './dto/create-transaction.dto';
import { UpdateTransactionDto } from './dto/update-transaction.dto';
import { CreateTransactionFromAnotherBankAccountDto } from './dto/create-transaction-from-another-bank-account.dto';
import { ConfirmTransactionDto } from './dto/confirm-transaction.dto';
import { InjectRepository } from '@nestjs/typeorm';
import { Transaction, TransactionOperation } from './entities/transaction.entity';
import { DataSource, Repository } from 'typeorm';
import { BankAccount } from 'src/bank-accounts/entities/bank-account.entity';
import { ClientKafka } from '@nestjs/microservices';
import { lastValueFrom } from 'rxjs';

@Injectable()
export class TransactionsService {
  constructor(
    @InjectRepository(Transaction)
    private transactionRepo: Repository<Transaction>,
    @InjectRepository(BankAccount)
    private bankAccountRepo: Repository<BankAccount>,
    private dataSource: DataSource,
    @Inject('KAFKA_SERVICE')
    private kafkaService: ClientKafka,
  ) {}

  async create(
    bankAccountId: string,
    createTransactionDto: CreateTransactionDto,
  ) {
    // DataSource.transaction refere-se a uma transação de banco de dados
    // A variável que recebe o retorno é uma entidade Transaction :P
    const transaction = await this.dataSource.transaction(async (manager) => {

      // Como vamos carregar o BankAccount para a memória antes de alterar o saldo,
      // é fundamental realizar um lock na linha.
      // Outras requisições sobre a mesma conta serão forçadas a esperar.
      const bankAccount = await manager.findOneOrFail(BankAccount, {
        where: { id: bankAccountId },
        lock: { mode: 'pessimistic_write' },
      });

      const transaction = manager.create(Transaction, {
        ...createTransactionDto,
        amount: createTransactionDto.amount * -1,
        bank_account_id: bankAccountId,
        operation: TransactionOperation.debit,
      });

      await manager.save(transaction);

      bankAccount.balance += transaction.amount;
      await manager.save(bankAccount);
      return transaction;
    });

    const sendData = {
      id: transaction.id,
      accountId: bankAccountId,
      amount: createTransactionDto.amount,
      pixKeyTo: createTransactionDto.pix_key_key,
      pixKeyKindTo: createTransactionDto.pix_key_kind,
      description: createTransactionDto.description,
    };

    // A API do serviço devolve um Observable
    // lastValueFrom transforma em Promise
    await lastValueFrom(this.kafkaService.emit('transactions', sendData));
    return transaction;
  }

  findAll(bankAccountId: string) {
    return this.transactionRepo.find({
      where: { bank_account_id: bankAccountId },
      order: { created_at: 'DESC' },
    });
  }

  findOne(id: number) {
    return `This action returns a #${id} transaction`;
  }

  update(id: number, updateTransactionDto: UpdateTransactionDto) {
    return `This action updates a #${id} transaction`;
  }

  remove(id: number) {
    return `This action removes a #${id} transaction`;
  }

  async confirmTransaction(arg0: ConfirmTransactionDto) {
    throw new Error('Method not implemented.');
  }

  async createFromAnotherBankAccount(
    input: CreateTransactionFromAnotherBankAccountDto
  ) {
    
  }
}
