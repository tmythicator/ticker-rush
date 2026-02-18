import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
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
      className={error ? 'border-destructive focus-visible:ring-destructive' : ''}
    />
    {error && <p className="text-xs text-destructive">{error}</p>}
  </div>
);
