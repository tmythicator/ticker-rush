import { FormField } from './FormField';
import { Input, type InputProps } from './Input';
import { type FieldValues, type Path, type UseFormRegister } from 'react-hook-form';

export interface FormInputProps<T extends FieldValues> extends Omit<InputProps, 'name'> {
  label: string;
  id: Path<T>;
  register: UseFormRegister<T>;
  error?: string;
}

export const FormInput = <T extends FieldValues>({
  label,
  id,
  register,
  error,
  ...props
}: FormInputProps<T>) => {
  return (
    <FormField label={label} htmlFor={id} error={error}>
      <Input id={id} {...register(id)} variant={error ? 'error' : 'default'} {...props} />
    </FormField>
  );
};
