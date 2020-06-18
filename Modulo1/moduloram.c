/*
 * Assignment 01 for the course Rootkit Programming at TUM in WS2014/15.
 * Implemented by Guru Chandrasekhara and Martin Herrmann
 */
#include <linux/proc_fs.h>
#include <linux/seq_file.h>
#include <asm/uaccess.h>
#include <linux/hugetlb.h>
#include <linux/module.h>
#include <linux/init.h>
#include <linux/kernel.h>
#include <linux/fs.h>

#define BUFSIZE 150

MODULE_LICENSE("GPL");
MODULE_DESCRIPTION("Escribe información de ram");
MODULE_AUTHOR("Berny Cardona - 201408603");

int sysinfo(struct sysinfo *info);

static int escribir_archivo(struct seq_file *m, void *v)
{
	struct sysinfo info;
	si_meminfo(&info);
	unsigned long memoria_libre = info.freeram * 4;
	unsigned long memoria_total = info.totalram * 4;
	unsigned long memoria_usada = memoria_total - memoria_libre;


	seq_printf(m,"Nombre Estudiante 1: Berny Andree Cardona Ramos\n");
	seq_printf(m,"Nombre Estudiante 2: Gary Stephen Giron Molina\n");
	seq_printf(m,"Carnet 1: 201408603\n");
	seq_printf(m,"Carnet 2: 201403997\n");
	seq_printf(m,"Memoria Total: %lu \n",memoria_total);
	seq_printf(m,"Memoria Libre: \t %lu \n",memoria_libre);
	
	seq_printf(m,"Porcentaje Utilizado: %li%% \n", (memoria_usada*100)/memoria_total);
	return 0;
}

static int al_abrir(struct inode *inode, struct file *file){
	return single_open(file, escribir_archivo, NULL);
}

static struct file_operations operaciones = 
{
	.open = al_abrir,
	.read = seq_read
};

static int iniciar(void)
{
	proc_create("memo_201408603", 0, NULL, &operaciones);
	printk(KERN_INFO "Carnet 1: 201408603\n Carnet 2: 201403997\n");
	return 0;
}

static void salir(void)
{
	remove_proc_entry("memo_201408603", NULL);
	printk(KERN_INFO "Vacaciones Junio: SOPES 1");
}

module_init(iniciar);
module_exit(salir);