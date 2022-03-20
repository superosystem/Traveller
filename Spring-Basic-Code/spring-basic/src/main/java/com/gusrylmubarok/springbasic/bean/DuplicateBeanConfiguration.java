package com.gusrylmubarok.springbasic.bean;

import com.gusrylmubarok.springbasic.data.Foo;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

/**
 * @author Gusryl Mubarok
 * @email gusrylmubarok@gmail.com
 */

@Configuration
public class DuplicateBeanConfiguration {

    @Bean
    public Foo foo1() {

        return new Foo();
    }

    @Bean
    public Foo foo2() {

        return new Foo();
    }
}
