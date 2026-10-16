import { AttachFile, DateRange, Delete, Label } from '@mui/icons-material';
import {
  Box,
  Button,
  Card,
  CardContent,
  IconButton,
  Popover,
  Stack,
  Typography,
} from '@mui/material';
import { useRef, useState } from 'react';
import Modal from '@/components/ui/Modal';
import Avatar from '@/components/ui/Avatar';
import { useNavigate } from 'react-router';
import datetime from '@/utils/datetime';
import common from '@/utils/common';

const ModalDetailTask = ({ open, handleClose, data }) => {
  const navigate = useNavigate();
  const [anchorElLabel, setAnchorElLabel] = useState(null);
  const [attachemnts, setAttachments] = useState([]);

  const handleClickLabel = (event) => {
    setAnchorElLabel(event.currentTarget);
  };

  const onClose = () => {
    setAnchorElLabel(null);
  };

  const openLabelPopover = Boolean(anchorElLabel);

  const refInputAttachments = useRef();

  const handleInputAttachments = (e) => {
    const files = e.target.files;
    if (files && files[0]) {
      console.log('files attachemnts: ', files);
      const currentAttachments = [...attachemnts];
      currentAttachments.push(files[0]);
      setAttachments(currentAttachments);
    }
  };

  const handleRemoveAttachments = (index) => {
    const currentAttachments = [...attachemnts];
    currentAttachments.splice(index, 1);
    setAttachments(currentAttachments);
  };

  return (
    <Modal
      open={open}
      handleClose={handleClose}
      aria-labelledby="modal-modal-title"
      aria-describedby="modal-modal-description"
      title={data?.title ?? 'Detail Tugas'}
      sx={{
        minWidth: 1300,
        maxWidth: 1300,
      }}
    >
      <Box
        sx={{
          display: 'flex',
          justifyContent: 'stretch',
          alignItems: 'stretch',
          overflowY: 'auto',
          p: 3,
          gap: 2,
        }}
      >
        <Box
          sx={{
            width: '60%',
            height: '100%',
          }}
        >
          <Box
            sx={{
              display: 'flex',
              justifyContent: 'flex-start',
              alignItems: 'center',
              gap: 1,
              mb: 3,
            }}
          >
            {/* <Button
              type="button"
              aria-describedby={'label-popover'}
              size="small"
              variant="outlined"
              startIcon={<Label />}
              onClick={handleClickLabel}
            >
              Label
            </Button> */}
            {/* <Popover
              id={'label-popover'}
              open={openLabelPopover}
              anchorEl={anchorElLabel}
              onClose={onClose}
              anchorOrigin={{
                vertical: 'bottom',
                horizontal: 'left',
              }}
            >
              <Typography sx={{ p: 2 }}>The content of the Popover.</Typography>
            </Popover> */}
            <Button size="small" variant="outlined" startIcon={<DateRange />}>
              {datetime.format(data?.due_date)}
            </Button>
          </Box>
          <Box
            sx={{
              mb: 3,
            }}
          >
            {/* <Typography variant="body1" fontWeight={'bold'} mb={2}>
              Pelanggan
            </Typography>
            <Box
              sx={{
                display: 'flex',
                gap: 1,
              }}
            >
              <Avatar
                label={'PO'}
                text={'Proyek Owner'}
                onClick={() => navigate('/customers')}
              />
            </Box> */}
          </Box>
          <Box>
            <Typography variant="body1" fontWeight={'bold'} mb={1}>
              Deskripsi
            </Typography>
            <Typography variant="body2">{data?.description}</Typography>
          </Box>
        </Box>
        <Box
          sx={{
            height: '100%',
            boxSizing: 'border-box',
            overflowY: 'auto',
            pb: 3,
          }}
        >
          <Typography variant="body1" fontWeight={'bold'} mb={2}>
            Lampiran (Gambar, Video, PDF)
          </Typography>
          <Typography variant="body2" mb={2}>
            Anda dapat menambahkan lampiran seperti gambar, video, atau PDF
            untuk memberikan informasi lebih lanjut tentang tugas ini.
          </Typography>
          <Box
            sx={{
              display: 'flex',
              flexDirection: 'column',
              gap: 1,
              mb: 2,
            }}
          >
            {attachemnts?.map((item, index) => (
              <Card key={index} variant="outlined">
                <CardContent>
                  <Stack
                    direction="row"
                    spacing={2}
                    alignItems="center"
                    justifyContent={'space-between'}
                  >
                    <Stack direction="row" spacing={1} alignItems="center">
                      <AttachFile />
                      <Box>
                        <Typography variant="body1">{item.name}</Typography>
                        <Typography variant="body2" color="textSecondary">
                          {common.bytesToMb(item.size)}
                        </Typography>
                      </Box>
                    </Stack>
                    <Box>
                      <IconButton
                        onClick={() => handleRemoveAttachments(index)}
                        color="error"
                        variant="outlined"
                        size="small"
                      >
                        <Delete />
                      </IconButton>
                    </Box>
                  </Stack>
                </CardContent>
              </Card>
            ))}
          </Box>
          <Button
            onClick={() => {
              refInputAttachments.current.click();
            }}
            variant="outlined"
            startIcon={<AttachFile />}
            fullWidth
          >
            Attachment
          </Button>
          <input
            style={{ display: 'none' }}
            type="file"
            ref={refInputAttachments}
            accept=".pdf, .mp4, .jpeg, .png, .jpg"
            onChange={handleInputAttachments}
          />
        </Box>
      </Box>
    </Modal>
  );
};

export default ModalDetailTask;
