import { NestFactory } from '@nestjs/core';
import { AppModule } from './app.module';
import { PixKeyAlreadyExistsFilter } from './pix-keys/filter/pix-key-already-exists.error';

async function bootstrap() {
  const app = await NestFactory.create(AppModule);
  app.useGlobalFilters(new PixKeyAlreadyExistsFilter());
  await app.listen(3000);
}
bootstrap();
