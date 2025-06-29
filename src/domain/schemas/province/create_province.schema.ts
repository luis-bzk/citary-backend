import { z } from 'zod';

export const CreateProvinceSchema = z.object({
  name: z
    .string({ required_error: 'El nombre de la provincia es requerido' })
    .min(1, 'El nombre de la provincia es requerido')
    .max(50, 'El nombre de la provincia no puede tener más de 50 caracteres')
    .transform((s) => s.toLowerCase()),
  code: z
    .string({ required_error: 'El código de la provincia es requerido' })
    .min(1, 'El código de la provincia es requerido')
    .max(10, 'El código de la provincia no puede tener más de 10 caracteres')
    .transform((s) => s.toLowerCase()),
  prefix: z
    .string({ required_error: 'El prefijo de la provincia es requerido' })
    .min(1, 'El prefijo de la provincia es requerido')
    .max(10, 'El prefijo de la provincia no puede tener más de 10 caracteres')
    .transform((s) => s.toLowerCase()),
  id_country: z
    .number({
      required_error: 'El ID del país es requerido',
      invalid_type_error: 'El ID del país debe ser un número válido',
    })
    .int('El ID del país debe ser un número entero')
    .min(1, 'El ID del país debe ser mayor que 0')
    .max(500, 'El ID del país debe ser menor o igual a 500'),
});

export type CreateProvinceDto = z.infer<typeof CreateProvinceSchema>;
