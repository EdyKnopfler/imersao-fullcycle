import { Controller, Get, Post, Body, Patch, Param, Delete, ValidationPipe } from '@nestjs/common';
import { TransactionsService } from './transactions.service';
import { CreateTransactionDto } from './dto/create-transaction.dto';
import { UpdateTransactionDto } from './dto/update-transaction.dto';
import { MessagePattern, Payload } from '@nestjs/microservices';
import { CreateTransactionFromAnotherBankAccountDto } from './dto/create-transaction-from-another-bank-account.dto';
import { ConfirmTransactionDto } from './dto/confirm-transaction.dto';

@Controller('transactions')
export class TransactionsController {
  constructor(private readonly transactionsService: TransactionsService) {}

  @Post()
  create(@Body() createTransactionDto: CreateTransactionDto) {
    return this.transactionsService.create(createTransactionDto);
  }

  @Get()
  findAll() {
    return this.transactionsService.findAll();
  }

  @Get(':id')
  findOne(@Param('id') id: string) {
    return this.transactionsService.findOne(+id);
  }

  @Patch(':id')
  update(@Param('id') id: string, @Body() updateTransactionDto: UpdateTransactionDto) {
    return this.transactionsService.update(+id, updateTransactionDto);
  }

  @Delete(':id')
  remove(@Param('id') id: string) {
    return this.transactionsService.remove(+id);
  }

  @MessagePattern(`bank${process.env.BANK_CODE}`)
  async onTransactionProcessedBank001(
    @Payload(new ValidationPipe())
    message: CreateTransactionFromAnotherBankAccountDto | ConfirmTransactionDto,
  ) {
    try {
      if (message.status === 'pending') {
        await this.transactionsService.createFromAnotherBankAccount(
          message as CreateTransactionFromAnotherBankAccountDto,
        );
      }
      if (message.status === 'confirmed') {
        await this.transactionsService.confirmTransaction(
          message as ConfirmTransactionDto,
        );
      }
    } catch (err) {
      console.error(err);
    }
  }

}
