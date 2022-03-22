package com.gusrylmubarok.hibernate;

import com.gusrylmubarok.hibernate.entity.Address;
import com.gusrylmubarok.hibernate.entity.Name;
import com.gusrylmubarok.hibernate.entity.User;
import com.gusrylmubarok.hibernate.util.HibernateUtil;
import org.hibernate.Session;
import org.hibernate.Transaction;

/**
 * Hello world!
 *
 */
public class App 
{
    public static void main( String[] args )
    {
        Name name = new Name("Agus", "Syahril", "Mubarok");
        Address address = new Address("111", "Ciawigebang", "Kuningan", "Jawa Barat", "Indonesia", "45591");
        User user = new User(name, "agus@gmail.com", address);

        Transaction transaction = null;
        try (Session session = HibernateUtil.getSessionFactory().openSession()){
            // Start
            transaction = session.beginTransaction();
            // save
            session.save(user);
            // commit
            transaction.commit();
        } catch (Exception ex) {
            if (transaction != null) {
                transaction.rollback();
            }
            ex.printStackTrace();
        }
    }
}
