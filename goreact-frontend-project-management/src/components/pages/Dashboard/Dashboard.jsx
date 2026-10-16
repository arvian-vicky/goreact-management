import { colors, Paper, Typography } from '@mui/material';

import SidebarLayout from '@/components/layouts/SidebarLayout';

const Dashboard = () => {
  return (
    <SidebarLayout pageTitle="Dashboard">
      <Paper
        sx={{
          padding: 2,
          background: colors.lightBlue[100],
        }}
      >
        <Typography>Menampilkan dashboard di sini</Typography>
      </Paper>
    </SidebarLayout>
  );
};

export default Dashboard;
