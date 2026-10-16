import { Box } from '@mui/material';
import { DatePicker as BaseDatePicker } from '@mui/x-date-pickers';
import { Controller } from 'react-hook-form';

const DatePicker = ({
  control,
  name,
  label,
  defaultValue,
  helperText,
  ...props
}) => {
  return (
    <Controller
      name={name}
      control={control}
      defaultValue={defaultValue}
      render={({ field: { value, onChange, onBlur } }) => {
        return (
          <Box
            sx={{
              marginBottom: 2,
            }}
          >
            <BaseDatePicker
              {...props}
              label={label}
              value={value}
              onChange={onChange}
              slotProps={{
                textField: {
                  fullWidth: true,
                  helperText,
                  onBlur,
                },
              }}
            />
          </Box>
        );
      }}
    />
  );
};

export default DatePicker;
