import TextField from '@/components/ui/Forms/TextField';
import { Box, Button, Stack } from '@mui/material';
import { useForm } from 'react-hook-form';
import Modal from '@/components/ui/Modal';
import Select from '@/components/ui/Forms/Select';
import cards from '@/services/api/cards';
import DatePicker from '@/components/ui/Forms/DatePicker';

const ModalAddNewTask = ({ open, handleClose, data, onSaved }) => {
  const {
    control,
    handleSubmit,
    formState: { errors },
    reset,
  } = useForm({
    defaultValues: {
      title: '',
      description: '',
      due_date: null,
    },
  });

  const onSubmit = async (payload) => {
    onSaved({
      list_id: data.id,
      title: payload.title,
      description: payload.description,
      due_date: payload.due_date,
    });

    reset();
    handleClose(); // Close the modal after submission
  };

  return (
    <Modal
      open={open}
      handleClose={handleClose}
      aria-labelledby="modal-modal-title"
      aria-describedby="modal-modal-description"
      title={'Tambah Tugas Baru'}
      sx={{
        minWidth: 400,
        maxWidth: 500,
      }}
    >
      <Box component={'form'} onSubmit={handleSubmit(onSubmit)}>
        <Stack p={2}>
          <TextField
            control={control}
            name="title"
            label="Judul Tugas"
            placeholder="Enter card title"
            fullWidth
          />
          <TextField
            control={control}
            name="description"
            label="Deskripsi"
            placeholder="Enter card description"
            multiline
            rows={4}
            helperText="Optional"
            fullWidth
          />
          <DatePicker
            control={control}
            label="Tanggal Deadline"
            name="due_date"
            sx={{
              width: '100%',
            }}
          />
        </Stack>
        <Stack
          direction={'row'}
          spacing={1}
          justifyContent="flex-end"
          sx={{ p: 2 }}
        >
          <Button type="submit" variant="contained">
            Simpan tugas
          </Button>
          <Button type="button" variant="outlined" onClick={handleClose}>
            Batal
          </Button>
        </Stack>
      </Box>
    </Modal>
  );
};

export default ModalAddNewTask;
