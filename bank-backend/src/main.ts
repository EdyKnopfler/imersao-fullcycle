import { NestFactory } from '@nestjs/core';
import { AppModule } from './app.module';
import { PixKeyAlreadyExistsFilter } from './pix-keys/filter/pix-key-already-exists.error';
import { ValidationPipe } from '@nestjs/common';

async function bootstrap() {
  const app = await NestFactory.create(AppModule);
  app.useGlobalFilters(new PixKeyAlreadyExistsFilter());
  app.useGlobalPipes(new ValidationPipe({ errorHttpStatusCode: 422 }));
  await app.listen(process.env.PORT || 3000);
}
bootstrap();
