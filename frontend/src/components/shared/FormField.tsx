import { Input } from '@/components/shared/Input';
import { Label } from '@/components/shared/Label';
import { type FieldValues, type Path, type UseFormRegister } from 'react-hook-form';

interface FormFieldProps<T extends FieldValues> {
  label: string;
  id: Path<T>;
  register: UseFormRegister<T>;
  error?: string;
  placeholder?: string;
  type?: string;
}

export const FormField = <T extends FieldValues>({
  label,
  id,
  register,
  error,
  placeholder,
  type = 'text',
}: FormFieldProps<T>) => (
  <div className="space-y-2">
    <Label htmlFor={id}>{label}</Label>
    <Input
      {...register(id)}
      id={id}
      type={type}
      placeholder={placeholder}
      variant={error ? 'error' : 'default'}
    />
    {error && <p className="text-xs text-destructive">{error}</p>}
  </div>
);
