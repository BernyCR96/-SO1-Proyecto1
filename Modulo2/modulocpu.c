
#include <linux/proc_fs.h>
#include <linux/seq_file.h>
#include <asm/uaccess.h>
#include <linux/hugetlb.h>
#include <linux/kernel.h>
#include <linux/module.h>
#include <linux/init.h>
#include <linux/sched/signal.h>
#include <linux/sched.h>
#include <linux/fs.h>
 
MODULE_LICENSE("GPL");
MODULE_DESCRIPTION("Escribe información de cpu");
MODULE_AUTHOR("Berny Cardona 201408603");
 
struct task_struct *task;
struct task_struct *task_child;
struct list_head *list;


static int escribir_archivo(struct seq_file *m, void *v)
{
    seq_printf(m,"Nombre Estudiante 1: Berny Cardona\n");
	seq_printf(m,"Nombre Estudiante 2: Gary Stephen Giron Molina\n");
	seq_printf(m,"Carnet 1: 201408603\n");
	seq_printf(m,"Carnet 2: 201403997\n");

    for_each_process( task ){
        seq_printf(m,"\nPid: %d Nombre: %s Estado: %ld\n",task->pid, task->comm,task->state);
        seq_printf(m,"Hijos:\n");
        list_for_each(list, &task->children)
        {
            task_child = list_entry(list, struct task_struct, sibling );
            seq_printf(m,"\nPid: %d Nombre: %s Estado: %ld\n",task_child->pid, task_child->comm,task_child->state);
        }
        seq_printf(m,"-----------------------------------------------------\n");
    }
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

static int simple_init(void)
{
	proc_create("cpu_201408603", 0, NULL, &operaciones);
	printk(KERN_INFO "Carnet 1: 201408603\n Carnet 2: 201403997\n");
	return 0;
}

static void simple_out(void)
{
	remove_proc_entry("cpu_201408603", NULL);
	printk(KERN_INFO "Vacaciones Junio: SOPES 1");
}
 
module_init(simple_init);
module_exit(simple_out);