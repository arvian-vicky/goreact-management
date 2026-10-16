import Box from '@mui/material/Box';
import Button from '@mui/material/Button';
import Typography from '@mui/material/Typography';

import { IconButton, Stack } from '@mui/material';
import Modal from '@/components/ui/Modal';
import { Close } from '@mui/icons-material';
import { useForm } from 'react-hook-form';
import DatePicker from '@/components/ui/Forms/DatePicker';
import boards from '@/services/api/boards';
import { useParams } from 'react-router';

const ModalUpdateDeadline = ({ open, data, handleClose }) => {
  const { id } = useParams();
  const { control, reset, handleSubmit } = useForm();

  const onSubmit = async (payload) => {
    await boards.update(id, {
      ...data,
      due_date: payload.due_date,
    });
    reset();
    handleClose();
  };

  return (
    <Modal
      open={open}
      handleClose={handleClose}
      aria-labelledby="modal-modal-title"
      aria-describedby="modal-modal-description"
      title={'Perbarui Deadline Proyek'}
      sx={{
        minWidth: 400,
        maxWidth: 500,
      }}
    >
      <Box component={'form'} onSubmit={handleSubmit(onSubmit)}>
        <Box
          sx={{
            p: 2,
          }}
        >
          <DatePicker
            control={control}
            label="Tanggal Deadline"
            name="due_date"
            sx={{
              width: '100%',
            }}
          />
        </Box>
        <Stack
          direction="row"
          spacing={1}
          p={2}
          justifyContent="flex-end"
          mt={3}
        >
          <Button type="submit" variant="contained" sx={{ mt: 2 }}>
            Simpan
          </Button>
          <Button
            variant="outlined"
            onClick={handleClose}
            sx={{ mt: 2, ml: 1 }}
            color="primary"
          >
            Batal
          </Button>
        </Stack>
      </Box>
    </Modal>
  );
};

export default ModalUpdateDeadline;
