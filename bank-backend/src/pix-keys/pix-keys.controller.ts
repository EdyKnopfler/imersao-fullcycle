import { Controller, Get, Post, Body, Patch, Param, Delete } from '@nestjs/common';
import { PixKeysService } from './pix-keys.service';
import { CreatePixKeyDto } from './dto/create-pix-key.dto';
import { UpdatePixKeyDto } from './dto/update-pix-key.dto';

@Controller('pix-keys')
export class PixKeysController {
  constructor(private readonly pixKeysService: PixKeysService) {}

  @Post()
  create(@Body() createPixKeyDto: CreatePixKeyDto) {
    return this.pixKeysService.create(createPixKeyDto);
  }

  @Get()
  findAll() {
    return this.pixKeysService.findAll();
  }

  @Get(':id')
  findOne(@Param('id') id: string) {
    return this.pixKeysService.findOne(+id);
  }

  @Patch(':id')
  update(@Param('id') id: string, @Body() updatePixKeyDto: UpdatePixKeyDto) {
    return this.pixKeysService.update(+id, updatePixKeyDto);
  }

  @Delete(':id')
  remove(@Param('id') id: string) {
    return this.pixKeysService.remove(+id);
  }
}
