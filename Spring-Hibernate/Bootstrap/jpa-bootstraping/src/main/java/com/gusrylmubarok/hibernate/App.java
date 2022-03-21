package com.gusrylmubarok.hibernate;

import com.gusrylmubarok.hibernate.entity.Student;
import com.gusrylmubarok.hibernate.util.JPAUtil;

import javax.persistence.EntityManager;

/**
 * Hello world!
 *
 */
public class App 
{
    public static void main( String[] args )
    {
        EntityManager entityManager = JPAUtil.getEntityManagerFactory().createEntityManager();
        entityManager.getTransaction().begin();

        Student student = new Student("Gusryl", "Mubarok", "gusrylmubarok@gmail.com");
        entityManager.persist(student);
        entityManager.getTransaction().commit();
        entityManager.close();

        JPAUtil.shutdown();
    }
}
